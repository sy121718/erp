# Application 应用层

> 目录：`internal/application/`

## 职责

- 业务流程编排
- 调用 Domain 层
- 处理事务
- **只操作 Entity，不涉及 DTO**

## Service 应用服务

```go
type AdminService struct {
    adminRepo domain.AdminRepository
}

// Create 创建管理员
func (s *AdminService) Create(ctx context.Context, req *CreateRequest) (*admin.Admin, error) {
    // 检查用户名是否存在
    exists, err := s.adminRepo.ExistsByUsername(ctx, req.Username)
    if err != nil {
        return nil, err
    }
    if exists {
        return nil, admin.ErrUsernameAlreadyExists
    }

    // 创建实体
    entity := &admin.Admin{
        Username: req.Username,
        Name:     req.Name,
        Email:    admin.Email(req.Email),
    }
    entity.SetPassword(req.Password)

    // 保存
    if err := s.adminRepo.Save(ctx, entity); err != nil {
        return nil, err
    }
    return entity, nil
}

// GetByID 根据ID获取
func (s *AdminService) GetByID(ctx context.Context, id int64) (*admin.Admin, error) {
    entity, err := s.adminRepo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }
    if entity == nil {
        return nil, admin.ErrAdminNotFound
    }
    return entity, nil
}
```

## 目录结构

```
internal/application/
└── admin/
    └── service.go    # 应用服务
```