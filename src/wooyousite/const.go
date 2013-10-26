/**
 * Title：
 *
 * Description:
 *
 * Author: Ivan
 *
 * Create Time: 2013-09-16 10:20
 *
 * Version: 1.0
 *
 * 修改历史: 版本号 修改日期 修改人 修改说明
 *   1.0 2013-09-16 Ivan 创建文件
 */
package wooyousite

//通用代码实体
type CodeEntity struct {
	Key   string `json:key`
	Value string `json:value`
}

// 表单的操作类型
type FormActionType struct {
	CREATE    CodeEntity `json:"create"`
	RETRIEVE  CodeEntity `json:"retrieve"`
	UPDATE    CodeEntity `json:"update"`
	DELETE    CodeEntity `json:"delete"`
	TOMBSTONE CodeEntity `json:"tombstone"`
}

// 表单的操作类型
var FORM_ACTION_TYPE FormActionType = FormActionType{
	CREATE:    CodeEntity{"C", "添加"},
	RETRIEVE:  CodeEntity{"R", "查询"},
	UPDATE:    CodeEntity{"U", "修改"},
	DELETE:    CodeEntity{"D", "删除"},
	TOMBSTONE: CodeEntity{"T", "逻辑删除"},
}

// 图片规格
const (
	IMAGE_SIZE_1 = "290x175"
	IMAGE_SIZE_2 = "290x370"
	IMAGE_SIZE_3 = "590x185"
	IMAGE_SIZE_4 = "600x370"
	IMAGE_SIZE_5 = "800x600"
)
