# Hướng dẫn tạo API CRUD (DDD + CQRS Architecture)

Tài liệu này là "Kim chỉ nam" để bạn xây dựng các tính năng mới cho dự án `platform` sao cho đồng bộ hoàn toàn với cấu trúc của dự án `store`.

---

## 1. Luồng Dữ Liệu (Data Flow)

Hệ thống hoạt động theo mô hình tách biệt luồng Đọc và Ghi:

*   **Luồng ĐỌC (Query)**: `Controller` ➔ `QueryService` ➔ `qsdto` (Dùng cho List, Fetch data thuần túy).
*   **Luồng GHI (Command)**: `Controller` ➔ `Usecase` ➔ `Repository` ➔ `dbmodel` (Dùng cho Create, Update, Delete).

---

## 2. Quy trình 8 bước triển khai chuẩn mực

### Bước 1: Định nghĩa API (OpenAPI First)
*   **Nơi thực hiện**: `src/openapi/openapi/`.
*   **Việc cần làm**: Định nghĩa Path, Request, Response và đặc biệt là `operationId`.
*   **Lệnh chạy**: `make codegen` (alias: `make gen`) — tái sinh mã trong `src/generated/`.

### Bước 2: Tạo Folder Domain riêng biệt
*   **Nơi thực hiện**: `src/internal/domain/<tên_nghiệp_vụ>/`.
*   **Quy tắc**: Mỗi module (user, coupon, order...) phải có folder riêng.
    *   **File nghiệp vụ**: `src/internal/domain/user/user.go`.
    *   **Nội dung**: Chứa Entity (Struct) và `Repository interface`. 
    *   **LƯU Ý**: Không tạo Interface cho Usecase tại đây.

### Bước 3: Định nghĩa Model Database (dbmodel)
*   **Nơi thực hiện**: `src/internal/infrastructure/dbmodel/`.
*   **Việc cần làm**: Khai báo struct mapping trực tiếp với bảng DB. Cần có hàm `ToDomain()` để chuyển đổi sang Entity ở Bước 2.

### Bước 4: Triển khai luồng ĐỌC (Query Service)
*   **DTO**: Tạo struct trong `src/internal/infrastructure/qsdto/`. Struct này chứa các tag `gorm` để query trực tiếp (ví dụ: `RankDto`).
*   **Service**: Tạo file trong `src/internal/infrastructure/queryservice/`. Viết hàm `Fetch...` để lấy dữ liệu đổ trực tiếp vào DTO.

### Bước 5: Triển khai luồng GHI (Repository)
*   **Nơi thực hiện**: `src/internal/infrastructure/repository/`.
*   **Việc cần làm**: Cài đặt (Implement) các interface đã khai báo ở Bước 2. Tập trung vào các hàm thay đổi dữ liệu (Insert, Update, Delete).

### Bước 6: Xử lý Logic nghiệp vụ (Usecase)
*   **Nơi thực hiện**: `src/internal/usecases/`.
*   **Cơ chế**: Dùng **Concrete Struct** (Struct cụ thể), không dùng interface.
*   **Hàm khởi tạo**: `func New...Usecase(repo ...) *...Usecase { ... }` (Trả về con trỏ struct).

### Bước 7: Điều khiển và Trả về kết quả (Controller)
*   **Nơi thực hiện**: `src/internal/interfaces/controller/`.
*   **Việc cần làm**:
    *   Tiêm (Inject) cả `Usecase` và `QueryService`.
    *   Sử dụng `common.BindJSON` để validate dữ liệu.
    *   Sử dụng `common.IsError...` để xử lý lỗi HTTP.

### Bước 8: Dependency Injection & Router
*   **Lệnh chạy**: `make wire`.
*   **Router**: Đăng ký controller mới vào router trong `src/cmd/server/main.go`.

---

## Quick examples & notes

- OperationId example (OpenAPI path):

    ```yaml
    post:
        operationId: createUser
        requestBody:
            content:
                application/json:
                    schema:
                        $ref: '../components/schemas/user.yml#/CreateUser'
        responses:
            '201':
                description: Created
    ```

- Generated code policy: `src/generated/` contains generated sources (models and server stubs). Regenerate with `make codegen`. Prefer committing generated code only when necessary for CI or downstream consumers.

- Migrations & seed: migrations live in `src/db/migrations/` and seed logic in `src/db/seed/`. Use your migration tooling or the included scripts to run migrations before starting the app.

- Transactions: perform transactional logic in repository or usecase layer; return errors to allow rollback. Keep transaction boundaries small and explicit.

- Tests & CI: run unit tests with `go test ./...`. For integration tests that need DB, use `testonly/` helpers and a test database container. Add CI step `make test` if desired.

- Validation & errors: use helpers in `internal/common` to unify error format and HTTP codes. Validate request bodies in controller (via generated models) and in usecase for business rules.

- Auth & permissions: implement authentication in `internal/middleware/` and enforce authorization in controller/usecase as appropriate.


---

## 3. Bản đồ Logic (Code để ở đâu?)

| Loại xử lý | Folder mục tiêu | Ví dụ |
| :--- | :--- | :--- |
| **Thiết kế API** | `openapi/` | `users.yml` |
| **Interface DB (Write)**| `domain/<module>/` | `UserRepository interface` |
| **GORM Model** | `infrastructure/dbmodel/` | `user.go` (DB Struct) |
| **Query Data (Read)** | `infrastructure/queryservice/` | `UserQueryService` |
| **Raw Result (DTO)** | `infrastructure/qsdto/` | `UserDto` |
| **Business Logic** | `usecases/` | `UserUsecase struct` |
| **HTTP Handling** | `interfaces/controller/` | `UserController` |

---

## 4. Các lưu ý quan trọng

1.  **Repository vẫn có hàm Get**: Chỉ dùng khi cần lấy dữ liệu để kiểm tra logic ngay trước khi Update/Delete. Nếu chỉ để hiện lên màn hình ➔ dùng `QueryService`.
2.  **Package Name**: Trong folder `domain/user`, package phải là `package user`. Tránh đặt tên package chung chung là `domain`.
3.  **Error Handling**: Luôn dùng các helper trong `internal/common` để đảm bảo format lỗi 400/500 đồng nhất toàn hệ thống.
