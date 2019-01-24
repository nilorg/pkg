package model

// JQueryDataTableResult jquery datatables result
// http://www.datatables.club/manual/server-side.html#returndata
type JQueryDataTableResult struct {
	Draw            int         `json:"draw"`            // 必要。上面提到了，Datatables发送的draw是多少那么服务器就返回多少。 这里注意，作者出于安全的考虑，强烈要求把这个转换为整形，即数字后再返回，而不是纯粹的接受然后返回，这是 为了防止跨站脚本（XSS）攻击。
	RecordsTotal    int64       `json:"recordsTotal"`    // 必要。即没有过滤的记录数（数据库里总共记录数）
	RecordsFiltered int64       `json:"recordsFiltered"` // 必要。过滤后的记录数（如果有接收到前台的过滤条件，则返回的是过滤后的记录数）
	Data            interface{} `json:"data"`            // 必要。表中中需要显示的数据。这是一个对象数组，也可以只是数组，区别在于 纯数组前台就不需要用 columns绑定数据，会自动按照顺序去显示 ，而对象数组则需要使用 columns绑定数据才能正常显示。
	Error           string      `json:"error"`           // 可选。你可以定义一个错误来描述服务器出了问题后的友好提示
}

// JQueryDataTablesParameters ...
type JQueryDataTablesParameters struct {
	Draw    int                       // 请求次数计数器
	Start   int                       // 第一条数据的起始位置
	Length  int                       // 每页显示的数据条数
	Columns []JQueryDataTablesColumns // 数据列
	Order   []JQueryDataTablesOrder   // 排序
	Search  JQueryDataTablesSearch    // 搜索
}

// JQueryDataTablesOrder 排序
type JQueryDataTablesOrder struct {
	Column int    // 排序的列的索引
	Dir    string // 排序模式
}

// JQueryDataTablesColumns 数据列
type JQueryDataTablesColumns struct {
	Data       string                 // Data 数据源
	Name       string                 // Name 名称
	Searchable bool                   // Searchable 是否可以被搜索
	Orderable  bool                   // Orderable 是否可以排序
	Search     JQueryDataTablesSearch // Search 搜索
}

// JQueryDataTablesSearch 搜索
type JQueryDataTablesSearch struct {
	Value string // Value 全局的搜索条件的值
	Regex bool   // Regex 是否为正则表达式处理
}
