package dto

// ======================
// Success Messages
// ======================
const (
	MsgSuccess       string = "สำเร็จ"
	MsgCreated       string = "สร้างข้อมูลสำเร็จ"
	MsgUpdated       string = "อัปเดตข้อมูลสำเร็จ"
	MsgDeleted       string = "ลบข้อมูลสำเร็จ"
	MsgFetched       string = "ดึงข้อมูลสำเร็จ"
	MsgOperationDone string = "ดำเนินการเสร็จสิ้น"
)

// ======================
// Error Messages (Generic)
// ======================
const (
	MsgError          string = "เกิดข้อผิดพลาด"
	MsgInvalidRequest string = "คำขอไม่ถูกต้อง"
	MsgNotFound       string = "ไม่พบข้อมูล"
	MsgPageNotFound   string = "ไม่พบหน้า"
	MsgConflict       string = "ข้อมูลซ้ำ"
	MsgUnauthorized   string = "ไม่ได้รับอนุญาต"
	MsgForbidden      string = "ไม่มีสิทธิ์เข้าถึง"
	MsgInternalError  string = "เกิดข้อผิดพลาดภายในระบบ"
	MsgServiceDown    string = "ไม่สามารถเชื่อมต่อบริการได้"
	MsgTimeout        string = "การดำเนินการหมดเวลา"
)

// ======================
// Auth / User Messages
// ======================
const (
	MsgLoginSuccess      string = "เข้าสู่ระบบสำเร็จ"
	MsgLoginFailed       string = "เข้าสู่ระบบไม่สำเร็จ"
	MsgLogoutSuccess     string = "ออกจากระบบสำเร็จ"
	MsgRegisterSuccess   string = "สมัครสมาชิกสำเร็จ"
	MsgRegisterFailed    string = "สมัครสมาชิกไม่สำเร็จ"
	MsgInvalidCredential string = "อีเมลหรือรหัสผ่านไม่ถูกต้อง"
	MsgTokenExpired      string = "โทเค็นหมดอายุ"
	MsgTokenInvalid      string = "โทเค็นไม่ถูกต้อง"
)

// ======================
// Validation Messages
// ======================
const (
	MsgValidationFailed   string = "ข้อมูลไม่ผ่านการตรวจสอบ"
	MsgRequiredField      string = "กรุณากรอกข้อมูลให้ครบถ้วน"
	MsgInvalidEmail       string = "อีเมลไม่ถูกต้อง"
	MsgPasswordTooShort   string = "รหัสผ่านสั้นเกินไป"
	MsgPasswordTooWeak    string = "รหัสผ่านไม่ปลอดภัย"
	MsgFieldExceedsLength string = "ข้อมูลเกินความยาวที่กำหนด"
)

// ======================
// Database / Storage
// ======================
const (
	MsgDBConnectionFailed string = "เชื่อมต่อฐานข้อมูลล้มเหลว"
	MsgDBQueryFailed      string = "การสืบค้นข้อมูลล้มเหลว"
	MsgDBInsertFailed     string = "เพิ่มข้อมูลไม่สำเร็จ"
	MsgDBUpdateFailed     string = "อัปเดตข้อมูลไม่สำเร็จ"
	MsgDBDeleteFailed     string = "ลบข้อมูลไม่สำเร็จ"

	MsgFileUploadSuccess string = "อัปโหลดไฟล์สำเร็จ"
	MsgFileUploadFailed  string = "อัปโหลดไฟล์ไม่สำเร็จ"
	MsgFileNotFound      string = "ไม่พบไฟล์"
	MsgFileDeleted       string = "ลบไฟล์สำเร็จ"
)

// ======================
// Network / External API
// ======================
const (
	MsgExternalAPIFailed string = "ไม่สามารถเชื่อมต่อ API ภายนอกได้"
	MsgExternalTimeout   string = "API ภายนอกไม่ตอบสนอง"
	MsgInvalidResponse   string = "ข้อมูลที่ได้รับไม่ถูกต้อง"
)

// ======================
// Permission / Role
// ======================
const (
	MsgRoleAssigned     string = "กำหนดสิทธิ์สำเร็จ"
	MsgRoleRevoked      string = "ยกเลิกสิทธิ์สำเร็จ"
	MsgPermissionDenied string = "สิทธิ์ไม่เพียงพอในการเข้าถึง"
)
