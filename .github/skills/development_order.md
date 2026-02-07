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
| 1.9 | `internal/core/model` | `base.go`<br>`query.go` | Completed | Base Struct: `ID`, `Created/UpdatedAt`, `DeletedAt`. **SaaS**: Add `TenantID` here for global isolation. | Low | [x] | [x] | 87.2% |
| 1.10| `pkg/utils/i18n` | `i18n.go` | Completed | **I18n**: Internationalization support (based on Viper/YAML). | Low | [x] | [x] | 100% |

---

## Phase 2: Core Business Identity (Services - Access & User)
*Goal: System identification, authorization, and user profiles.*

| ID | Module Path | Feature / File | Status | Description | Diff | Test | Pass | Cov |
|:---|:---|:---|:---|:---|:---|:---|:---|:---|
| 2.1 | `internal/services/access` | `token/jwt.go` | Completed | **Token**: JWT generation, parsing, validation, refresh token logic. | Med | [x] | [x] | 87% |
| 2.2 | `internal/services/access` | `permission/casbin.go` | Completed | **Permission**: RBAC enforcement using Casbin or similar. Policy loader. | High | [x] | [x] | 83% |
| 2.3 | `internal/services/access` | `verify/otp.go` | Completed | **Verify**: Logic for generating/checking OTP codes (Email/SMS). | Med | [x] | [x] | 90% |
| 2.4 | `internal/services/access` | `operation/audit.go` | Completed | **Operation**: Structure for recording critical user operations (Audit Logs). | Low | [x] | [x] | 100% |
| 2.5 | `internal/services/user` | `account/auth.go` | Completed | **Account**: Register (Email/Phone), Login (Pwd), Logout logic. | High | [x] | [x] | 83% |
| 2.6 | `internal/services/user` | `account/password.go` | Completed | **Account**: Password hashing (bcrypt/argon2), change pwd, reset pwd. | Med | [x] | [x] | 100% |
| 2.7 | `internal/services/user` | `info/profile.go` | Completed | **Info**: Profile CRUD (Avatar, Bio, Gender, Birthday). | Low | [x] | [x] | 95% |
| 2.8 | `internal/services/user` | `preference/settings.go`| Completed | **Preference**: User specific configs (Theme, Notif settings). JSON storage? | Low | [x] | [x] | 92% |
| 2.9 | `internal/services/user` | `group/role.go` | Completed | **Group**: User groups/roles association (Admin, Editor, Member). | Med | [x] | [x] | 93% |

---

## Phase 3: Content & Interaction
*Goal: Content management and user interaction.*

| ID | Module Path | Feature / File | Status | Description | Diff | Test | Pass | Cov |
|:---|:---|:---|:---|:---|:---|:---|:---|:---|
| 3.1 | `internal/services/contents` | `category.go` | Completed | Taxonomy tree structure (Parent/Child). | Med | [x] | [x] | 90% |
| 3.2 | `internal/services/contents` | `article.go`<br>`banner.go` | Completed | Standard content CRUD. Publishing status workflow. | Med | [x] | [x] | 94% |
| 3.3 | `internal/services/contents` | `media.go` | Completed | Media library metadata (Links to OSS/Local files). | Low | [x] | [x] | 97% |
| 3.4 | `internal/services/shieldword`| `filter.go` | Completed | Text censoring/filtering logic (Sensitive word replacement). | Med | [x] | [x] | 92% |
| 3.5 | `internal/services/relation` | `follow.go` | Completed | User relationships (Follow/Fan). | Low | [x] | [x] | 88% |
| 3.6 | `internal/services/message` | `notification.go` | Completed | In-app notification creation and reading status. | Med | [x] | [x] | 95% |
| 3.7 | `internal/services/form` | `schema.go`<br>`submission.go` | Completed | **Custom Forms**: JSON Schema definition and generic submission handling. | High | [x] | [x] | 91% |

---

## Phase 4: Commerce & Finance
*Goal: Transactions, Assets, and Orders.*

| ID | Module Path | Feature / File | Status | Description | Diff | Test | Pass | Cov |
|:---|:---|:---|:---|:---|:---|:---|:---|:---|
| 4.1 | `internal/services/world` | `saas/tenant.go` | Completed | **SaaS**: Tenant existence, domain resolution, and configuration. | Med | [x] | [x] | 87% |
| 4.2 | `internal/services/commerce`| `product/sku.go` | Completed | Product SPUs (Display) and SKUs (Stock keeping units). | High | [x] | [x] | 86% |
| 4.3 | `internal/services/commerce`| `stock/inventory.go` | Completed | Inventory management (Deduct/Restore mechanisms). | High | [x] | [x] | 70% |
| 4.4 | `internal/services/commerce`| `coupon/rule.go` | Completed | Coupon distribution and validity rules. | Med | [x] | [x] | 76% |
| 4.5 | `internal/services/commerce`| `writeoff/verify.go` | Completed | **Verification**: QR Code/Code verification logic for offline usage. | Med | [x] | [x] | 100% |
| 4.6 | `internal/services/commerce`| `order/fsm.go` | Completed | Order State Machine (Created -> Paid -> Shipped -> Completed). | High | [x] | [x] | 100% |
| 4.7 | `internal/services/user` | `pocket/assets.go` | Completed | **User Wallet**: Aggregated view of Points/Coupons (No Cash top-up). | Med | [x] | [x] | 100% |
| 4.8 | `internal/services/finance` | `ledger.go` | Completed | **Finance**: Transaction logs (`Deals`) and Point history. | Med | [x] | [x] | 91.4% |

---

## Phase 5: Adapters & Integrations
*Goal: Connect to external systems.*

| ID | Module Path | Feature / File | Status | Description | Diff | Test | Pass | Cov |
|:---|:---|:---|:---|:---|:---|:---|:---|:---|
| 5.1 | `pkg/extra/cloudstorage` | `uploader.go` | Completed | Cloud implementations (AliyunOSS/AWS S3) for `SimpleUploader`. | Med | [x] | [x] | 90% |
| 5.2 | `pkg/extra/sms` | `sender.go` | Completed | SMS Provider integration (Aliyun/Twilio). | Low | [x] | [x] | 100% |
| 5.3 | `pkg/thirdparty` | `payment/pay.go` | Completed | Unified Payment Interface. Adapters for WeChat/Stripe. | High | [x] | [x] | 100% |

---

## Phase 6: Presentation (APIS)
*Goal: Public/User-facing JSON APIs.*

| ID | Module Path | Feature / File | Status | Description | Diff | Test | Pass | Cov |
|:---|:---|:---|:---|:---|:---|:---|:---|:---|
| 6.1 | `internal/core/route` | `engine.go`<br>`middleware.go` | Completed | Gin Framework setup, SaaS Middleware, Logger Middleware. | Med | [x] | [x] | 80% |
| 6.2 | `internal/apis/auth` | `handler.go` | Completed | **Auth**: Login (Pwd/Bypass), Register, Change Password. | Low | [x] | [x] | 90% |
| 6.3 | `internal/apis/account`| `handler.go` | Completed | **Account**: Update Profile, Get Detail, List. | Low | [x] | [x] | 85% |
| 6.4 | `internal/apis/content`| `handler.go` | Completed | **Content**: Article/Banner CRUD. | Med | [x] | [x] | 80% |
| 6.5 | `internal/apis/redirect`| `handler.go` | Completed | **Redirect**: WeChat Login, OSS Callback. | Med | [x] | [x] | 80% |

## Phase 7: Presentation (Admin)
*Goal: Admin Management Panel APIs.*

| ID | Module Path | Feature / File | Status | Description | Diff | Test | Pass | Cov |
|:---|:---|:---|:---|:---|:---|:---|:---|:---|
| 7.1 | `internal/admin/auth` | `handler.go` | Completed | **Admin Auth**: Login. | Low | [x] | [x] | 80% |
| 7.2 | `internal/admin/user` | `handler.go` | Completed | **User**: Manage Users (List, Detail, Ban). | Med | [x] | [x] | 80% |
| 7.3 | `internal/admin/system` | `handler.go` | Completed | **System**: Menu Config & Init. | Low | [x] | [x] | 80% |
| 7.4 | `internal/admin/ui` | `shell.go` | Completed | **Frontend**: SPA Shell (React/Mantine/CDN) + Login. | Med | [x] | [x] | 0% |
| 7.5 | `web/admin/src` | `App.jsx` | Completed | **Frontend**: Layout & Dynamic Menu (from JSON) & UserList. | High | [x] | [x] | 0% |
| 7.6 | `internal/admin/contents`| `handler.go` | Completed | **Content**: Manage Articles, Categories. | Med | [x] | [x] | 80% |
| 7.7 | `internal/admin/commerce`| `handler.go` | Pending | **Commerce**: Product, Order, Coupons. | High | [ ] | [ ] | 0% |

---

## Phase 8: Assembly
*Goal: Application Entry.*

| ID | Module Path | Feature / File | Status | Description | Diff | Test | Pass | Cov |
|:---|:---|:---|:---|:---|:---|:---|:---|:---|
| 8.1 | `cmd/appsite-monolith` | `main.go` | Completed | Dependency Injection (Wire). Start HTTP Server. | Med | [x] | [x] | 0% |

---

## Technical Debt & TODOs
*Goal: Track remaining tasks marked in code.*

| Type | Module Path | File | Description | Status |
|:---|:---|:---|:---|:---|
| TODO | `internal/admin` | `router.go` | Add Admin Middleware (Route Protection). | Pending |
| TODO | `internal/apis/redirect` | `handler.go` | Implement WeChat OAuth callback logic. | Pending |
| TODO | `internal/apis/redirect` | `handler.go` | Implement OSS Callback logic. | Pending |
| TODO | `internal/apis/content` | `handler.go` | Assign AuthorID from context (Auth). | Pending |
