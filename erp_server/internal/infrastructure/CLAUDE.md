# Infrastructure 基础设施层

> 目录：`internal/infrastructure/`

## 职责

- Repository 实现
- PO（持久化对象）定义
- PO ↔ Entity 转换
- 外部服务调用

## PO 持久化对象

```go
type AdminPO struct {
    ID         int64          `gorm:"primaryKey"`
    Username   string         `gorm:"column:username;uniqueIndex"`
    Password   string         `gorm:"column:password"`
    Name       string         `gorm:"column:name"`
    Email      string         `gorm:"column:email"`
    Status     int            `gorm:"column:status"`
    CreateTime time.Time      `gorm:"column:create_time"`
    UpdateTime time.Time      `gorm:"column:update_time"`
    DeletedAt  gorm.DeletedAt `gorm:"column:deleted_time;index"`
}

func (AdminPO) TableName() string {
    return "sys_admin"
}
```

## Converter 转换器

```go
// ToEntity PO → Entity
func ToEntity(po *AdminPO) *admin.Admin {
    return &admin.Admin{
        ID:        po.ID,
        Username:  po.Username,
        Password:  po.Password,
        Name:      po.Name,
        Email:     admin.Email(po.Email),
        Status:    admin.Status(po.Status),
        CreateTime: po.CreateTime,
    }
}

// ToPO Entity → PO
func ToPO(entity *admin.Admin) *AdminPO {
    return &AdminPO{
        ID:       entity.ID,
        Username: entity.Username,
        Password: entity.Password,
        Name:     entity.Name,
        Email:    string(entity.Email),
        Status:   int(entity.Status),
    }
}
```

## Repository 实现

```go
type AdminRepository struct {
    db *gorm.DB
}

func NewAdminRepository() admin.AdminRepository {
    return &AdminRepository{db: database.Get()}
}

func (r *AdminRepository) FindByID(ctx context.Context, id int64) (*admin.Admin, error) {
    var po AdminPO
    err := r.db.WithContext(ctx).Where("id = ?", id).First(&po).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    return ToEntity(&po), nil
}

func (r *AdminRepository) Save(ctx context.Context, entity *admin.Admin) error {
    po := ToPO(entity)
    return r.db.WithContext(ctx).Create(&po).Error
}
```

## 目录结构

```
internal/infrastructure/
└── repository/
    └── admin/
        ├── repository.go  # 仓储实现
        └── po.go          # 持久化对象 + 转换器
```