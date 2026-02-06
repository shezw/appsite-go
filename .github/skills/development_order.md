# Progressive Development Plan & Status

This document tracks the development progress of the Appsite Go framework.
**Requirement**: All modules generally require unit tests. Checkboxes `[ ]` track progress.

## Table Legend
*   **Status**: `Pending` (待开发), `In Progress` (开发中), `Testing` (测试中), `Fixing` (修复中), `Completed` (已完成)
*   **Diff**: Difficulty (Low/Med/High)
*   **Test**: `_test.go` exists?
*   **Pass**: Tests strictly passed?
*   **Cov**: Coverage % (Aim 100%, Min 70%; Single Fun > 50%)

---

## Phase 1: Foundation (Core & Utils)
*Goal: Establish robust infrastructure and stateless utilities.*

| ID | Module Path | Feature / File | Status | Description | Diff | Test | Pass | Cov |
|:---|:---|:---|:---|:---|:---|:---|:---|:---|
| 1.1 | `internal/core/error` | `types.go`<br>`code.go` | Completed | Define `AppError` struct, interfaces, and standardized error codes (e.g., HTTP map). | Low | [x] | [x] | 79% |
| 1.2 | `internal/core/log` | `logger.go`<br>`zap.go` | Completed | Structured logging wrapper (Zap implementation recommended). Context-aware logging. | Low | [x] | [x] | 91% |
| 1.3 | `pkg/utils/file` | `file.go`<br>`check.go` | Completed | File existence check, mkdir, size formatting, mime-type detection. | Low | [x] | [x] | 83% |
| 1.4 | `pkg/utils/timeconvert` | `time.go` | Completed | Date formatting, parsing, duration helpers, timezone utilities. | Low | [x] | [x] | 92% |
| 1.5 | `pkg/utils/simpleimage` | `resize.go`<br>`check.go` | Completed | Basic image resizing, format checking (Validation before upload). | Med | [x] | [x] | 77% |
| 1.6 | `internal/core/setting` | `loader.go`<br>`config.go` | Completed | Configuration struct definitions. Load from YAML/Env. Hot-reload support. | Med | [x] | [x] | 91.7% |
| 1.7 | `pkg/utils/orm` | `gorm_init.go`<br>`scopes.go` | Completed | GORM (or Ent) initialization config, common scopes (Pagination, SoftDelete). | Med | [x] | [x] | 96.7% |
| 1.8 | `pkg/utils/redis` | `client.go`<br>`lock.go` | Completed | Redis client init. Distributed lock helper implementation. | Med | [x] | [x] | 93.8% |
| 1.9 | `internal/core/model` | `base.go`<br>`query.go` | Pending | Base Struct: `ID`, `Created/UpdatedAt`, `DeletedAt`. **SaaS**: Add `TenantID` here for global isolation. | Low | [ ] | [ ] | 0% |

---

## Phase 2: Core Business Identity (Services - Access & User)
*Goal: System identification, authorization, and user profiles.*

| ID | Module Path | Feature / File | Status | Description | Diff | Test | Pass | Cov |
|:---|:---|:---|:---|:---|:---|:---|:---|:---|
| 2.1 | `internal/services/access` | `token/jwt.go` | Pending | **Token**: JWT generation, parsing, validation, refresh token logic. | Med | [ ] | [ ] | 0% |
| 2.2 | `internal/services/access` | `permission/casbin.go` | Pending | **Permission**: RBAC enforcement using Casbin or similar. Policy loader. | High | [ ] | [ ] | 0% |
| 2.3 | `internal/services/access` | `verify/otp.go` | Pending | **Verify**: Logic for generating/checking OTP codes (Email/SMS). | Med | [ ] | [ ] | 0% |
| 2.4 | `internal/services/access` | `operation/audit.go` | Pending | **Operation**: Structure for recording critical user operations (Audit Logs). | Low | [ ] | [ ] | 0% |
| 2.5 | `internal/services/user` | `account/auth.go` | Pending | **Account**: Register (Email/Phone), Login (Pwd), Logout logic. | High | [ ] | [ ] | 0% |
| 2.6 | `internal/services/user` | `account/password.go` | Pending | **Account**: Password hashing (bcrypt/argon2), change pwd, reset pwd. | Med | [ ] | [ ] | 0% |
| 2.7 | `internal/services/user` | `info/profile.go` | Pending | **Info**: Profile CRUD (Avatar, Bio, Gender, Birthday). | Low | [ ] | [ ] | 0% |
| 2.8 | `internal/services/user` | `preference/settings.go`| Pending | **Preference**: User specific configs (Theme, Notif settings). JSON storage? | Low | [ ] | [ ] | 0% |
| 2.9 | `internal/services/user` | `group/role.go` | Pending | **Group**: User groups/roles association (Admin, Editor, Member). | Med | [ ] | [ ] | 0% |

---

## Phase 3: Content & Interaction
*Goal: Content management and user interaction.*

| ID | Module Path | Feature / File | Status | Description | Diff | Test | Pass | Cov |
|:---|:---|:---|:---|:---|:---|:---|:---|:---|
| 3.1 | `internal/services/contents` | `category.go` | Pending | Taxonomy tree structure (Parent/Child). | Med | [ ] | [ ] | 0% |
| 3.2 | `internal/services/contents` | `article.go`<br>`banner.go` | Pending | Standard content CRUD. Publishing status workflow. | Med | [ ] | [ ] | 0% |
| 3.3 | `internal/services/contents` | `media.go` | Pending | Media library metadata (Links to OSS/Local files). | Low | [ ] | [ ] | 0% |
| 3.4 | `internal/services/shieldword`| `filter.go` | Pending | Text censoring/filtering logic (Sensitive word replacement). | Med | [ ] | [ ] | 0% |
| 3.5 | `internal/services/relation` | `follow.go` | Pending | User relationships (Follow/Fan). | Low | [ ] | [ ] | 0% |
| 3.6 | `internal/services/message` | `notification.go` | Pending | In-app notification creation and reading status. | Med | [ ] | [ ] | 0% |
| 3.7 | `internal/services/form` | `schema.go`<br>`submission.go` | Pending | **Custom Forms**: JSON Schema definition and generic submission handling. | High | [ ] | [ ] | 0% |

---

## Phase 4: Commerce & Finance
*Goal: Transactions, Assets, and Orders.*

| ID | Module Path | Feature / File | Status | Description | Diff | Test | Pass | Cov |
|:---|:---|:---|:---|:---|:---|:---|:---|:---|
| 4.1 | `internal/services/world` | `saas/tenant.go` | Pending | **SaaS**: Tenant existence, domain resolution, and configuration. | Med | [ ] | [ ] | 0% |
| 4.2 | `internal/services/commerce`| `product/sku.go` | Pending | Product SPUs (Display) and SKUs (Stock keeping units). | High | [ ] | [ ] | 0% |
| 4.3 | `internal/services/commerce`| `stock/inventory.go` | Pending | Inventory management (Deduct/Restore mechanisms). | High | [ ] | [ ] | 0% |
| 4.4 | `internal/services/commerce`| `coupon/rule.go` | Pending | Coupon distribution and validity rules. | Med | [ ] | [ ] | 0% |
| 4.5 | `internal/services/commerce`| `writeoff/verify.go` | Pending | **Verification**: QR Code/Code verification logic for offline usage. | Med | [ ] | [ ] | 0% |
| 4.6 | `internal/services/commerce`| `order/fsm.go` | Pending | Order State Machine (Created -> Paid -> Shipped -> Completed). | High | [ ] | [ ] | 0% |
| 4.7 | `internal/services/user` | `pocket/assets.go` | Pending | **User Wallet**: Aggregated view of Points/Coupons (No Cash top-up). | Med | [ ] | [ ] | 0% |
| 4.8 | `internal/services/finance` | `ledger.go` | Pending | **Finance**: Transaction logs (`Deals`) and Point history. | Med | [ ] | [ ] | 0% |

---

## Phase 5: Adapters & Integrations
*Goal: Connect to external systems.*

| ID | Module Path | Feature / File | Status | Description | Diff | Test | Pass | Cov |
|:---|:---|:---|:---|:---|:---|:---|:---|:---|
| 5.1 | `pkg/extra/cloudstorage` | `uploader.go` | Pending | Cloud implementations (AliyunOSS/AWS S3) for `SimpleUploader`. | Med | [ ] | [ ] | 0% |
| 5.2 | `pkg/extra/sms` | `sender.go` | Pending | SMS Provider integration (Aliyun/Twilio). | Low | [ ] | [ ] | 0% |
| 5.3 | `pkg/thirdparty` | `payment/pay.go` | Pending | Unified Payment Interface. Adapters for WeChat/Stripe. | High | [ ] | [ ] | 0% |

---

## Phase 6: Presentation (HTTP)
*Goal: Routing and JSON APIs.*

| ID | Module Path | Feature / File | Status | Description | Diff | Test | Pass | Cov |
|:---|:---|:---|:---|:---|:---|:---|:---|:---|
| 6.1 | `internal/core/route` | `engine.go`<br>`middleware.go` | Pending | Gin/Favorite framework setup. **SaaS Middleware** (Tenant Extraction). | Med | [ ] | [ ] | 0% |
| 6.2 | `internal/apis` | `*/handler.go` | Pending | Map HTTP Requests -> Service Interfaces -> HTTP Responses. | Low | [ ] | [ ] | 0% |
| 6.3 | `internal/admin` | `*/handler.go` | Pending | Admin panel specific endpoints (Requires Admin Auth). | Low | [ ] | [ ] | 0% |

---

## Phase 7: Assembly
*Goal: Application Entry.*

| ID | Module Path | Feature / File | Status | Description | Diff | Test | Pass | Cov |
|:---|:---|:---|:---|:---|:---|:---|:---|:---|
| 7.1 | `cmd/appsite-monolith` | `main.go` | Pending | Dependency Injection (Wire). Start HTTP Server. | Med | [ ] | [ ] | 0% |
