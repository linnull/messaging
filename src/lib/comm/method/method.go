package method

// 大类
const method_type_system = 1
const method_type_logic = 2

// 系统小类
const system_command = 1

// 逻辑小类
const logic_login = 1

// 方法
const MethodId_system_command = (method_type_system << 24) | (system_command << 8) | 1 // 0x1000101

const MethodId_login = (method_type_logic << 24) | (logic_login << 8) | 1 // 0x2000101

// 对应方法名
