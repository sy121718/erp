# Domain 领域层

> 目录：`internal/domain/`

## 职责

- 核心业务规则
- 实体定义（Entity）
- 值对象定义
- 仓储接口定义

## Entity 实体

```go
type Admin struct {
    ID        int64
    Username  string
    Password  string
    Name      string
    Email     Email    // 使用值对象
    Status    Status
    CreatedAt time.Time
}

// SetPassword 设置密码（业务方法）
func (a *Admin) SetPassword(password string) error {
    if len(password) < 6 {
        return ErrPasswordTooShort
    }
    a.Password = hashPassword(password)
    return nil
}
```

## Value Object 值对象

```go
type Email string

func NewEmail(s string) (Email, error) {
    if !isValidEmail(s) {
        return "", ErrInvalidEmail
    }
    return Email(s), nil
}

type Status int

const (
    StatusDisabled Status = 0
    StatusEnabled  Status = 1
)
```

## Repository 仓储接口

```go
type AdminRepository interface {
    FindByID(ctx context.Context, id int64) (*Admin, error)
    FindByUsername(ctx context.Context, username string) (*Admin, error)
    ExistsByUsername(ctx context.Context, username string) (bool, error)
    Save(ctx context.Context, admin *Admin) error
    Update(ctx context.Context, admin *Admin) error
    Delete(ctx context.Context, id int64) error
}
```

## 领域错误

```go
var (
    ErrAdminNotFound         = errors.NewNotFound(404001, "管理员不存在")
    ErrUsernameAlreadyExists = errors.NewBadRequest(10001, "用户名已存在")
    ErrInvalidPassword       = errors.NewBadRequest(10002, "密码错误")
)
```

## 目录结构

```
internal/domain/
├── admin/
│   ├── admin.go           # 实体
│   ├── value_objects.go   # 值对象
│   ├── repository.go      # 仓储接口
│   └── errors.go          # 领域错误
└── common/                # 公共领域组件
```