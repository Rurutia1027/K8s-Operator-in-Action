# Learning Makefile — run from repo root: make -C learning <target>

REPO_ROOT := $(abspath ..)
export REPO_ROOT

.PHONY: help test-all test-all-no-aws setup-envtest

help:
	@echo "Learning targets:"
	@echo "  make test-all-no-aws     Run issues 1-11 (no AWS)"
	@echo "  make test-all            Run all runnable Go tests"
	@echo "  make test-issue-NN       Run a single issue (01-11)"
	@echo "  make setup-envtest       Download envtest binaries"

setup-envtest:
	$(MAKE) -C $(REPO_ROOT) setup-envtest

# --- M1 ---
test-issue-01:
	cd $(REPO_ROOT) && go test ./learning/M1-api-and-crd/issue-01-crd-api-types/... -v

test-issue-02:
	./M1-api-and-crd/issue-02-install-crd/run.sh

# --- M2 ---
test-issue-03:
	cd $(REPO_ROOT) && go test ./learning/M2-controller-basics/issue-03-minimal-reconcile/... -v

test-issue-04:
	cd $(REPO_ROOT) && go test ./learning/M2-controller-basics/issue-04-finalizer/... -v

test-issue-05:
	cd $(REPO_ROOT) && go test ./learning/M2-controller-basics/issue-05-status-updates/... -v

# --- M3 ---
test-issue-06:
	cd $(REPO_ROOT) && go test ./learning/M3-fake-aws-client/issue-06-ec2-client-interface/... -v

test-issue-07:
	cd $(REPO_ROOT) && go test ./learning/M3-fake-aws-client/issue-07-full-reconcile-fake/... -v

test-issue-08:
	cd $(REPO_ROOT) && go test ./learning/M3-fake-aws-client/issue-08-envtest-suite/... -v

# --- M4 ---
test-issue-09:
	./M4-cluster-deploy/issue-09-kind-deploy/run.sh

test-issue-10:
	cd $(REPO_ROOT) && go test ./learning/M4-cluster-deploy/issue-10-drift-detection/... -v

test-issue-11:
	cd $(REPO_ROOT) && go test ./learning/M4-cluster-deploy/issue-11-spec-to-run-input/... -v

# --- M5 (manual, needs AWS) ---
test-issue-12:
	./M5-real-aws/issue-12-create-ec2/run.sh

test-issue-13:
	./M5-real-aws/issue-13-delete-ec2/run.sh

test-issue-14:
	./M5-real-aws/issue-14-e2e-validation/run.sh

test-all-no-aws: test-issue-01 test-issue-06 test-issue-07 test-issue-08 test-issue-10 test-issue-11 setup-envtest
	cd $(REPO_ROOT) && go test ./learning/M2-controller-basics/... -v

test-all: setup-envtest test-all-no-aws
	cd $(REPO_ROOT) && go test ./learning/... -v
