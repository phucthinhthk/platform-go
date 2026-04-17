# Hướng Dẫn Cấu Trúc Source Code & Quy Trình Phát Triển 

Cấu trúc source code trong dự án của bạn đang được tổ chức theo chuẩn **Clean Architecture** (Kiến trúc Sạch) kết hợp **Domain-Driven Design (DDD)** nhưng có một số tùy chỉnh (customization) rất đặc trưng. Mọi thứ được phân lớp nghiêm ngặt để đảm bảo code dễ bảo trì.

Dưới đây là một bản phân tích chi tiết tận răng xem code nằm ở đâu, hoạt động thế nào và quy trình lập trình cho cấu trúc gốc của bạn.

---

## 1. Sơ Đồ Cấu Trúc Mã Nguồn

```text
my-project/
├── src/                          # THƯ MỤC SOURCE CHÍNH CỦA DỰ ÁN
│   ├── cmd/                      # Nơi chứa file main.go chạy khởi động server
│   │   └── server/               
│   ├── config/                   # Nơi tải cấu hình (.env, .yaml) vào các biến structs của Go
│   ├── db/                       # Chứa các script init/migration db
│   │   ├── migrations/           
│   │   └── seed/                 
│   ├── generated/                # KHU VỰC CODE TỰ ĐỘNG SINH KHI DÙNG OPENAPI
│   │   ├── api/                  # Chứa code sau khi quét file YAML
│   │   │   ├── models.gen.go     
│   │   │   └── server.gen.go     
│   │   └── openapi/              # Các cấu hình trung gian Openapi nếu cần sinh tiếp
│   ├── internal/                 # THÙNG CHỨA TRANH LOGIC & BUSINESS CỦA ỨNG DỤNG
│   │   ├── common/               # Xử lý logic dùng chung (helpers, error code, types base)
│   │   ├── constants/            # Các biến/hằng số dùng chung toàn system
│   │   ├── controller/           # Chứa các file Route HTTP Handler. Nhận Request -> Call Logic -> Response
│   │   ├── controller_test/      # File unit test cho lớp controller
│   │   ├── di/                   # Chứa cấu hình Dependency Injection (Google Wire)
│   │   ├── domain/               # Trái tim lõi nhất: Chỉ chứa Interface và Entity của Domain, không chứa logic code
│   │   ├── infrastructure/       # Lớp ngoài cùng của Clean Architecture: Các kết nối tới "Thế Giới Bên Ngoài"
│   │   │   ├── auth/             # Kết nối server SSO/Auth
│   │   │   ├── dbmodel/          # GORM Model. Trực tiếp định nghĩa các struct bảng database (cột, type)
│   │   │   ├── dbutil/           # Utilities cho db như kết nối, context
│   │   │   ├── qsdto/            # Query Service DTO (CQRS Pattern)
│   │   │   ├── queryservice/     # Implement tầng Query Service cho CQS/CQRS (Read only)
│   │   │   └── repository/       # Implement các interface query Database cho Domain (Create, Update, Delete)
│   │   ├── interfaces/           # Lớp kết nối HTTP/gRPC API bên ngoài (có thể chứa custom HTTP handler mapper)
│   │   ├── middleware/           # HTTP Middleware xử lý trước khi vào Controller
│   │   ├── pkg/                  # Các công cụ helper nội bộ
│   │   └── usecases/             # BỘ NÃO CHÍNH: Nơi chứa toàn bộ cốt truyện và Logic Nghiệp vụ
│   ├── logs/                     # Thư mục lưu log file khi chạy ứng dụng
│   ├── mocks/                    # Chứa mocked interface/struct cho Unit Test
│   ├── openapi/                  # ĐỊNH NGHĨA API OPENAPI GỐC
│   │   ├── components/           # Khai báo Schema, Schema Params
│   │   ├── openapi/              # Nơi chứa file gốc định nghĩa OpenAPI
│   │   └── paths/                # Định nghĩa từng route con ở đây
│   ├── testonly/                 # Các tiện ích (utilities) hoặc data mockup test
│   └── tests/                    # Nơi run các file Integration / E2E test cho hệ thống
├── .air.toml                     # Cấu hình auto-reload code bằng Air khi gõ
└── .env.local / .env.test        # Lưu biến môi trường cục bộ/test (KHÔNG ĐẨY LÊN GIT)
```

---

## 2. Giải Đáp Các Câu Hỏi Chuyên Sâu

### A. API được cấu hình ở đâu trên OpenAPI? Viết file YAML ở đâu để sinh ra code?
* **Nơi khai báo:** Thiết kế API theo chuẩn OpenAPI đã được tách ra nhiều cấu phần nằm tại `src/openapi/` (`components`, `paths`, file gộp `openapi`). Việc bóc tách giúp dễ dàng quản lý.
* **Nơi sinh ra code:** Sau khi bạn viết xong định nghĩa route ở trong các file YAML (ví dụ bạn thêm route `GET /orders` trong `paths`). Bạn sẽ chạy công cụ sinh code (có định cấu hình bằng `make generate`), code Golang tạo ra sẽ trút xuống thư mục `src/generated/api/models.gen.go` và `server.gen.go`.

### B. Logic được xử lý chỗ nào? (Luồng Chạy Data Flow)
Luồng chạy của hệ thống này sẽ đi 1 chiều duy nhất từ ngoài vào trong lõi như sau:
`Request` ➔ `Middleware` ➔ `Controller` ➔ `Usecase` (hoặc `Query Service`) ➔ `Infrastructure`

1. **`src/internal/middleware/`**: Logic đánh chặn người dùng. Thấy chữ *JWT* thì xử lý xác thực đăng nhập trước.
2. **`src/internal/controller/`**: Logic bóc mổ. Tầng Controller trích xuất thông tin HTTP Request Body (đã được sinh bởi `models.gen.go`) gửi vào. Không viết code xử lý nghiệp vụ hay check Data SQL ở đây! Controller gọi hàm `Usecase` hoặc `Query Service`.
3. **`src/internal/usecases/`**: **LOGIC KINH DOANH CHÍNH NẰM Ở ĐÂY**.
   * Ví dụ: Usecase *"Tạo Đơn Hàng"* sẽ kiểm tra giá tiền hợp lệ, tính toán phí giao hàng, trạng thái kho, sau đó nó sẽ mượn tầng infrastructure để lưu xuống. 
4. **`src/internal/infrastructure/`**: Logic liên lạc với "bên ngoài thế giới Golang app":
   * Hệ thống áp dụng **CQRS / CQS Pattern**:
     * `repository/`: Phụ trách các lệnh Command (Create, Update, Delete) và chạy qua cấu trúc `domain` sử dụng `dbmodel/`.
     * `queryservice/` & `qsdto/`: Phụ trách Data Query thao tác lấy dữ liệu thuần túy (Read) và trả về DTO luôn (để bớt thao tác convert Domain model ở phía Usecase). 
   * `mail/`, `aws/`: Gửi Email cho khách hay tải file lên S3 vân vân được code tại đây.

### C. Vì sao phải bóc tách `dbmodel` khỏi `domain`?
* Trong source code này, thiết kế đã cẩn thận tách `dbmodel` (Tầng hạ tầng data layer trực tiếp gỡ cột SQL) và `domain` (Tầng nghiệp vụ thuần túy không quan tâm đến GORM/MySQL). 
* Usecase sẽ chỉ mượn interface của Layer `domain`. Layer `repository` sẽ cài đặt các interface domain đó nhưng ở bên dưới nó xài struct `dbmodel` để móc data từ Database lên, sau đó nó Convert sang struct thuần của `domain` rổi gửi trả về về `usecase`.

### D. File `di/wire.go` hoạt động ra sao?
Cấu trúc có chia layer rất tuyệt nhưng làm sao để khởi tạo chúng? Không nhẽ hàm main.go phải khai báo bằng tay dài dằng dặc `c := NewController(NewUsecase(NewRepository(), NewDBConn(), NewAWSS3()))`?.
Với Google Wire (`wire.go`), bạn chỉ cần khai báo Provider (ai cung cấp module nào). Lúc build hay chuẩn bị test, bạn chạy lệnh CLI `wire`, thư viện sẽ tự đọc code, nối dây và sinh ra file `wire_gen.go` xử lý hết mọi khai báo tiêm phụ thuộc.

---

## 3. Quy Trình (Rule) Triển Khai Tính Năng Chuẩn Mực Trong Source Này

Giả sử bạn phải làm tính năng: **Lấy danh sách các đơn hàng (Dealers Orders List)**.
Hãy tuân thủ nghiêm ngặt 5 bước này:

1. **API First (Config API Contract)**
   * Khai báo mới trong `src/openapi/paths/` cho URL path, và `src/openapi/components/` cho model.
   * Cập nhật file gốc (nếu dùng chung) và chạy công cụ sinh mã `make generate`. Code Controller Interfaces và Struct Req/Res sẽ sinh trong `src/generated/api/`.
   
2. **Setup Interfaces (Domain / Query Service)**
   * Vì đây là một Query (Lấy Data), tốt nhất hãy định nghĩa 1 Query Service (`OrderQueryService`) cùng với `OrderDTO` trong `src/internal/infrastructure/qsdto/`.
   * Nếu thao tác Command (Update/Delete/Create), bạn quy định interface trong `src/internal/domain/order.go`.
   
3. **Setup Logic Giao Tiếp DB (Infrastructure Layer)**
   * Định nghĩa cấu trúc bảng nếu chưa có trong `src/internal/infrastructure/dbmodel/order_model.go`.
   * Với Query: Lập trình logic cho `OrderQueryService` trong `src/internal/infrastructure/queryservice/`.
   * Với Command: Lập trình hàm `Save()`, `Delete()` trong `src/internal/infrastructure/repository/`.
   * Cập nhật các hàm khởi tạo trong Wire dependency `src/internal/di/wire.go`.
   
4. **Viết Cốt Truyện Kinh Doanh (Usecase Layer / Query Handler Layer)**
   * Nếu phức tạp, dùng Usecase (`src/internal/usecases/`). Nhưng nếu là Read đơn giản, Tầng Query Service có thể đẩy dữ liệu trực tiếp lên cho Controller hoặc qua 1 lớp Handler trung gian.

5. **Giao Tiếp Khách Ngoại Cấp (Controller Layer)**
   * Tạo hoặc mở file ở thư mục `src/internal/controller/`.
   * Lập trình struct Controller để implemement interface `server.gen.go` vừa đẻ ra ở Bước 1. (Bạn có thể viết unit test riêng trong `src/internal/controller_test/`).
   * Hàm ở trong Controller sẽ lấy param -> gọi `QueryService`/`Usecase` -> Lấy mảng kết quả -> map sang array các `models.gen.go` model -> Trả về qua HTTP. Hoàn thiện tính năng!
