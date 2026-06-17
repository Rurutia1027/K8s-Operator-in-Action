# 从 0 到 1：给 Spring Boot 单体服务做一个 Kubernetes 自定义 Operator（单 Controller 完整版）

这篇文章基于你当前仓库的学习脉络（API/CRD -> Reconcile -> Finalizer -> Status -> 集群部署），但把业务对象换成一个非常简单的 Spring Boot 单体服务。  
目标是：**只用 1 个 Controller**，把 Spring 服务在生产里最关键的运维能力都纳入 Operator。

---

## 1. 我们要解决什么问题

假设你有一个最小 Spring Boot 服务（例如 `demo-spring-app`），你希望做到：

- 通过一个自定义资源声明服务规模、镜像、配置、灰度参数
- Operator 自动创建/更新 Deployment 和 Service
- 删除自定义资源时，执行清理逻辑（Finalizer）
- 状态可观测：Ready、副本数、当前镜像、最近错误
- 支持常见能力：滚动升级、配置变更触发重启、暂停发布、回滚入口

一句话：把“写 YAML + 手工运维”的流程，变成“声明意图 + Controller 持续收敛”。

---

## 2. 最终架构（单 Controller）

我们定义一个 CRD：`SpringApp`，由一个 `SpringAppReconciler` 管理：

- **输入（Spec）**：镜像、副本数、端口、环境变量、ConfigMap 引用、发布策略
- **输出（Status）**：ObservedGeneration、Ready 条件、AvailableReplicas、Phase、LastError
- **关联资源**：Deployment、Service（可选 HPA/Ingress，本文先留接口）

单 Controller 的责任：

1. 读取 `SpringApp`
2. 对比期望状态与实际状态（Deployment/Service）
3. 创建或更新下游资源
4. 回写状态
5. 处理删除路径（Finalizer）

---

## 3. 准备工程（Kubebuilder）

> 如果你已有 operator 工程，可以跳过初始化，只做 API 和 Controller。

```bash
# 1) 初始化（示例）
kubebuilder init --domain example.com --repo github.com/your-org/spring-app-operator

# 2) 创建 API + Controller
kubebuilder create api --group apps --version v1alpha1 --kind SpringApp --controller --resource
```

目录通常会出现：

- `api/v1alpha1/springapp_types.go`
- `internal/controller/springapp_controller.go`
- `config/crd/`
- `config/rbac/`
- `config/manager/`

---

## 4. 设计 CRD（关键字段一个不少）

下面是建议的 `SpringAppSpec` / `SpringAppStatus` 结构（核心字段）：

```go
type SpringAppSpec struct {
    Image           string            `json:"image"`
    Replicas        *int32            `json:"replicas,omitempty"`
    ContainerPort   int32             `json:"containerPort,omitempty"`
    ServicePort     int32             `json:"servicePort,omitempty"`
    Env             map[string]string `json:"env,omitempty"`
    ConfigMapName   string            `json:"configMapName,omitempty"`
    Resources       corev1.ResourceRequirements `json:"resources,omitempty"`

    // 发布控制
    Paused          bool              `json:"paused,omitempty"`
    MaxUnavailable  *intstr.IntOrString `json:"maxUnavailable,omitempty"`
    MaxSurge        *intstr.IntOrString `json:"maxSurge,omitempty"`
}

type SpringAppStatus struct {
    ObservedGeneration int64              `json:"observedGeneration,omitempty"`
    Phase              string             `json:"phase,omitempty"` // Pending/Progressing/Ready/Failed
    AvailableReplicas  int32              `json:"availableReplicas,omitempty"`
    ReadyReplicas      int32              `json:"readyReplicas,omitempty"`
    CurrentImage       string             `json:"currentImage,omitempty"`
    LastError          string             `json:"lastError,omitempty"`
    Conditions         []metav1.Condition `json:"conditions,omitempty"`
}
```

建议加上校验注解（最小值、必填、枚举等），避免非法配置进入 Reconcile。

完成 API 修改后，执行：

```bash
make generate
make manifests
```

这一步会生成 DeepCopy 和 CRD YAML。  
**不要手改 `zz_generated.deepcopy.go`。**

---

## 5. RBAC 权限（经常漏）

Controller 至少需要以下权限：

- `springapps`、`springapps/status`、`springapps/finalizers`
- `deployments`（get/list/watch/create/update/patch/delete）
- `services`（get/list/watch/create/update/patch/delete）
- `events`（create/patch，可选但推荐）

在 `controller` 文件顶部通过 kubebuilder 注解声明，执行 `make manifests` 自动更新 RBAC。

---

## 6. Reconcile 主流程（单 Controller 模板）

一个可落地的顺序如下：

1. **Get SpringApp**  
   - NotFound 直接返回（对象已被删）
2. **处理删除分支**  
   - 若有 `DeletionTimestamp`，执行 `reconcileDelete`
3. **确保 Finalizer**  
   - 首次看到对象时添加 finalizer 并更新
4. **暂停发布检查**  
   - `spec.paused=true` 时只更新状态，不推进 Deployment 变更
5. **对齐 Service**  
   - 创建/更新 Service（selector 固定）
6. **对齐 Deployment**  
   - 组装期望 Deployment（镜像、资源、滚动策略、env、config）
   - Create 或 Patch（推荐 Server-Side Apply 或 controllerutil + mutate）
7. **聚合状态并回写 Status**  
   - `AvailableReplicas`、`Ready` 条件、`ObservedGeneration`
8. **重试与回队列策略**  
   - 出错返回 `error`，由控制器框架指数退避重试

---

## 7. Finalizer 删除路径（生产必须有）

删除 `SpringApp` 时，至少做这三件事：

1. 标记状态 `Phase=Deleting`（可选）
2. 执行外部清理（本文示例没有云资源，可留扩展点）
3. 移除 finalizer，让对象真正删除

即使当前只有 Deployment/Service，建议保留 finalizer 框架。  
后续你接入 DNS、证书、云资源（例如 EC2、EIP）时，不用重写删除语义。

---

## 8. 状态与可观测性（Status 不是装饰）

至少维护这些状态：

- `ObservedGeneration`：标记 controller 已处理到哪一代 spec
- `Ready` condition：True/False/Unknown + reason/message
- `AvailableReplicas`、`ReadyReplicas`
- `LastError`：最近一次错误摘要（避免直接塞大段堆栈）

推荐状态流转：

- 新建：`Pending`
- 创建中：`Progressing`
- 就绪：`Ready`
- 连续失败：`Failed`

结合 `kubectl get springapps -o wide`，你会得到很清晰的运维可见性。

---

## 9. Spring Boot 相关的“必须纳入 Operator”的功能

下面这些能力建议都由同一个 Controller 管理：

### 9.1 镜像升级与滚动策略

- spec 改 `image` -> Deployment template hash 变化 -> 滚动升级
- 支持 `maxUnavailable/maxSurge`，避免全量抖动

### 9.2 配置变更自动重启

- `configMapName` 内容变化时，需要触发 Pod 重建
- 常见做法：把 ConfigMap resourceVersion 或 checksum 注入 Pod annotation

### 9.3 弹性伸缩入口

- 先支持 `spec.replicas`
- 未来可扩展 `spec.autoscaling`，由同一 controller 负责 HPA 子资源

### 9.4 暂停发布

- `spec.paused=true` 时阻止 Deployment 模板变更下发
- 适合故障窗口冻结

### 9.5 回滚入口（最简版）

- 在 spec 里保留 `rollbackToRevision`（可选）
- controller 根据 revision 注解恢复历史模板（进阶）

---

## 10. 本地验证闭环（不依赖 AWS 也能完整练习）

你之前卡在 AWS 计费，这个路径可以完全绕开：

### 10.1 单元测试（优先）

- 构造 fake client，验证：
  - 首次 Reconcile 会创建 Deployment/Service
  - 修改 spec 会更新 Deployment
  - 删除路径会移除 finalizer
  - paused 模式不推进 rollout

### 10.2 envtest（推荐）

- 启动 apiserver/etcd 的测试环境
- 加载 CRD
- 跑真实 Reconcile 与 status 更新断言

### 10.3 Kind 集群验证

```bash
kind create cluster --name spring-operator
make install
make run
kubectl apply -f config/samples/apps_v1alpha1_springapp.yaml
kubectl get springapps -w
kubectl get deploy,svc,pod
```

### 10.4 观察点

- 改镜像 tag：看 rollout
- 改 ConfigMap：看是否触发重建
- 删 CR：看 finalizer 与资源清理

---

## 11. 生产化建议（你后续可以按这个顺序加）

1. **准入校验**：webhook 拦截非法配置  
2. **高可用**：controller leader election  
3. **指标与告警**：暴露 reconcile error rate、duration  
4. **事件审计**：关键动作写 K8s Event  
5. **版本演进**：`v1alpha1 -> v1beta1` conversion 策略  
6. **灰度能力**：canary 字段 + service selector 分流（后续仍可保持单 controller）

---

## 12. 常见坑（你现在这个阶段最容易踩）

- 忘了 `status` 子资源权限，导致状态更新报错
- Finalizer 逻辑里 return 太早，对象删不掉
- 直接改 informer cache 对象，而不是 DeepCopy 后再 mutate
- Deployment selector / labels 不一致，导致升级异常
- Reconcile 每次都 update，触发无意义写放大和事件风暴

---

## 13. 一份最小样例 CR（可直接改造）

```yaml
apiVersion: apps.example.com/v1alpha1
kind: SpringApp
metadata:
  name: demo-spring-app
  namespace: default
spec:
  image: ghcr.io/your-org/demo-spring-app:1.0.0
  replicas: 2
  containerPort: 8080
  servicePort: 80
  configMapName: demo-spring-app-config
  paused: false
```

---

## 14. 结语

你现在完全可以先把“Spring 单体服务 Operator 化”做扎实，再回头接 AWS。  
这样做有两个直接收益：

- **工程闭环先跑通**：不被云账号、账单、权限阻塞
- **控制面思维先建立**：声明式 API、Reconcile、Finalizer、Status 这些核心能力与云厂商无关

当你以后把 SpringApp 扩展到“创建/绑定云资源”（例如 EC2、EBS、LB）时，本质只是把外部依赖接到既有 Reconcile 框架中，而不是推倒重来。

如果你愿意，我下一步可以继续给你补一版：  
**“直接对照你当前仓库结构的落地清单（每个文件该改什么）”**，包括 `api/v1`、`internal/controller`、`config/samples`、测试用例骨架，一次性可开工。
