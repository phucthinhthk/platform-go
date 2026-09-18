# Thông số kỹ thuật xác thực và thu thập thông tin xác thực

## Tổng quan

Phần giải thích sau đây mô tả về phía quản lý cửa hàng, nhưng phía ứng dụng gốc nhìn chung là tương tự, vì vậy vui lòng điều chỉnh cho phù hợp.

Hệ thống này là hệ thống quản lý thông tin (tương lai sẽ phát triển nhiều mảng), và người dùng phải đăng nhập. Người dùng đăng nhập bằng cách cung cấp địa chỉ email và mật khẩu họ muốn đăng nhập (hiện tại mới chỉ có quyền admin, nên chỉ có admin mới login được, tương lai sẽ phát triển them). Thông tin xác thực được thu thập bằng mã thông báo JWT, và thông tin người dùng được lưu trữ trong phiên (ngữ cảnh Gin).

## Phương thức xác thực

### JWT (JSON Web Token)

- **Mã thông báo truy cập**: Có hiệu lực trong 24 giờ, được sử dụng để xác thực API

- **Mã thông báo làm mới**: Có hiệu lực trong 7 ngày, được sử dụng để cấp lại mã thông báo truy cập
- **Mã thông báo CSRF**: Có hiệu lực trong 24 giờ, ngăn chặn các cuộc tấn công CSRF

### Quản lý Cookie
Tất cả các mã thông báo được quản lý dưới dạng cookie chỉ HTTP:

| Tên Cookie | Mục đích | Ngày hết hạn | Đường dẫn | Chỉ HTTP | Bảo mật |

|----------|------|----------|------|----------|--------|

| `jwt` | Mã truy cập | 24 giờ | `/` | ✓ | ✓ |

| `refresh_token` | Mã làm mới | 7 ngày | `/refresh` | ✓ | ✓ |

| `csrf` | Mã CSRF | 24 giờ | `/` | ✓ | ✗ |

## Quy trình đăng nhập

### 1. Yêu cầu đăng nhập
```http

POST /login
Content-Type: application/json
{
"email": "admin@gmail.com",

"password": "123456",

}
```

### 2. Quy trình xác thực
1. Xác minh địa chỉ email và mật khẩu
2. Kiểm tra quyền truy cập cửa hàng của người dùng
3. Tạo mã thông báo (truy cập, làm mới, CSRF)
4. Thiết lập cookie

### 3. Kiểm tra quyền
Các kiểm tra quyền khác nhau được thực hiện tùy thuộc vào loại tài khoản người dùng:

- **Admin**: Có thể truy cập tất cả các trang
- **Người dùng**: Có thể một vài trang nhất định (hiện tại chưa phát triển)

## Xác thực Middleware

### Middleware JWTAuth
Xác thực JWT được thực hiện cho tất cả các điểm cuối API (ngoại trừ `/health` và `/login`).

```go
// Các đường dẫn cần bỏ qua
skipPaths := []string{"/login"}
```

### Quy trình xác thực
1. Lấy mã thông báo JWT từ cookie
2. Xác minh chữ ký mã thông báo và ngày hết hạn
3. Lấy thông tin người dùng từ các yêu cầu JWT
4. Lấy thông tin chi tiết người dùng từ cơ sở dữ liệu
5. Thiết lập thông tin trong ngữ cảnh

## Quản lý phiên

### Khóa ngữ cảnh
Sau khi xác thực, các thông tin sau được lưu trữ trong ngữ cảnh Gin:

```go
// Hằng số khóa ngữ cảnh
const (
LoginAdministratorKey = "login_administrator" // Thông tin người dùng đăng nhập
)
```

### Cấu trúc dữ liệu

#### LoginAdministrator
```go
type LoginAdministrator struct {
ID uint // ID người dùng
Name string // Tên người dùng
Email string // Địa chỉ email
AccountType string // Loại tài khoản
}
```
## Cách sử dụng
### Truy xuất thông tin trong Controller

#### Truy xuất thông tin người dùng
```go
func (controller) SomeAction(c *gin.Context) {
/ Lấy thông tin người dùng
loginUser, ok := c.Get(ctxkeys.LoginAdministratorKey)
if !ok {
c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi máy chủ nội bộ"})
return
}

/ Kiểm tra kiểu dữ liệu
user := loginUser.(types.LoginAdministrator)
userID := user.ID
userName := user.Name
accountType := user.AccountType

/ Ví dụ kiểm tra quyền // phần này chưa cần làm, tạm thời bỏ qua
if user.AccountType == "staff" { c.JSON(http.StatusForbidden, gin.H{"error": "Quyền truy cập bị từ chối"})
return
}

// Được sử dụng trong quá trình xử lý tiếp theo
result, err := someUseCase.Execute(c.Request.Context(), userID)
// ...

```

## Đăng xuất
### Yêu cầu đăng xuất
```http
POST /logout
```
### Quy trình đăng xuất
1. Xóa tất cả cookie (đặt ngày hết hạn thành quá khứ)
2. Trả về phản hồi thành công
```go
// Xóa Cookie
c.SetCookie("jwt", "", -1, "/", "", false, true)
c.SetCookie("csrf", "", -1, "/", "", false, false)
c.SetCookie("refresh_token", "", -1, "/refresh", "", (sai, đúng)
```
## Làm mới Token
### Yêu cầu làm mới
```http
POST /refresh
```
### Quy trình làm mới
1. Lấy token làm mới từ cookie
2. Xác minh token làm mới
3. Tạo token truy cập và token làm mới mới
4. Đặt token mới vào cookie
## Truy xuất thông tin người dùng
### Truy xuất thông tin đăng nhập hiện tại
```http
GET /me
```
Trả về thông tin người dùng được truy xuất từ ​​ngữ cảnh dưới dạng phản hồi.

## Các vấn đề về bảo mật

### Bảo vệ Token
- Tất cả token được quản lý bằng cookie chỉ dành cho HTTP.

- Cờ bảo mật chỉ được bật trong môi trường HTTPS.

- Token CSRF ngăn chặn tấn công giả mạo yêu cầu liên trang.

### Mật khẩu
- Được mã hóa và lưu trữ trong cơ sở dữ liệu.

- Được xác thực bằng giá trị mã hóa trong quá trình đăng nhập.

### Kiểm soát truy cập
- Bao gồm ID cửa hàng và ID dự án trong các yêu cầu JWT.

- Kiểm soát truy cập chi tiết dựa trên loại tài khoản.

- Xác thực lại thông tin người dùng trong cơ sở dữ liệu.

### Xử lý lỗi
- Không tiết lộ thông tin chi tiết trong trường hợp lỗi xác thực.

- Phản hồi lỗi thống nhất.

## Các tệp triển khai

### Các tệp chính
- **Bộ điều khiển xác thực**: `src/internal/controller/auth_controller.go`
- **Sử dụng đăng nhập**

Trường hợp**: `src/internal/usecases/auth/usecase/login_usecase.go`
- **Middleware JWT**: `src/internal/middleware/jwt_middleware.go`
- **Triển khai JWT**: `src/internal/infrastructure/security/jwt.go`
- **Định nghĩa kiểu**: `src/internal/pkg/types/login_user.go`
- **Khóa ngữ cảnh**: `src/internal/common/ctxkeys/context_key.go`
- **Định nghĩa OpenAPI**: `src/openapi/paths/auth.yml`

### Các tập tin kiểm thử
- **Kiểm thử đăng nhập**: `src/internal/controller_test/auth_controller/login_test.go`

- **Kiểm thử JWT**: `src/internal/infrastructure/security/jwt_test.go`

## Tóm tắt

Chức năng đăng nhập của hệ thống này sử dụng xác thực dựa trên thông tin người dùng, cho phép mỗi người dùng đăng nhập vào hệ thông để thực hiện các tác vụ của họ. Bảo mật được đảm bảo thông qua xác thực JWT, và việc quản lý phiên được thực hiện bằng cách sử dụng chức năng ngữ cảnh của Gin. Bộ điều khiển có thể dễ dàng truy xuất thông tin người dùng bằng cách sử dụng `ctxkeys.LoginShopKey` và `ctxkeys.LoginAdministratorKey`.