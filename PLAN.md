# Appsite Monolith Development Plan

Based on the architecture design, here is the implementation order for the modules.

## Phase 1: Core Foundation (Completed)
- [x] Project Structure
- [x] Configuration (Viper)
- [x] Logger (Zap)
- [x] Database (GORM, SQLite/MySQL)
- [x] Redis (Go-Redis, Miniredis)
- [x] I18n Utility

## Phase 2: APIS - Account & Authentication (In Progress)
- [x] User Entity (GORM)
- [x] Auth Services (Login, Register, Token)
- [x] API Handlers (Login, Register)
- [ ] Account Operations
    - [ ] `loginByPass` (Bypass/OTP Login)
    - [ ] `changePass` (Change Password)
    - [ ] `update` (Update Profile)

## Phase 3: APIS - Content Management
- [ ] Content Service
- [ ] APIs
    - [ ] `add`
    - [ ] `update`
    - [ ] `delete`
    - [ ] `list`
    - [ ] `detail`

## Phase 4: APIS - Redirect / Social
- [ ] `wechatLogin`
- [ ] `ossCallback`

## Phase 5: Admin - Core & User
- [ ] `install` (System Installation)
- [ ] `dashboard` (Statistics)
- [ ] `setting` (System Settings)
- [ ] User Management
    - [ ] `login` (Admin Login)
    - [ ] `regist` (Admin Create User)
    - [ ] `account` (List/Manage)
    - [ ] `profile`

## Phase 6: Admin - Contents
- [ ] `article`
- [ ] `category`

## Phase 7: Admin - Commerce
- [ ] `product`
- [ ] `coupon`
- [ ] `order`
- [ ] `payment`
- [ ] `shipping`

## Phase 8: Admin - System
- [ ] `config`
- [ ] `database`
- [ ] `performance`
- [ ] `FontEnd` (Admin UI Assets)
