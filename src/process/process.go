package process
import (
	"fmt"
        "github.com/hjqhezgh/lessgo" 
)

type node struct {
        id int
	process_type int
	back_id int
	pass_id int
	node_type string
	role string
}

type order struct {
        id int
	process_type int
	current_node int
	create_user int
	create_time string
}

type task struct {
	id int 
	order_id int
	node_id int
	actor int
	action string
	result string
}

type Process struct {
	order *order
	task_id int
	process_type int
	current_node *node
	nodes []*node
	tasks []*task
}

func Test(){
	var process = new(Process)
	fmt.Println("process testing")
	//process = CreateProcess(2)
	process = GetProcess(87)
	fmt.Println(process.order)
	fmt.Println(process.current_node)
	for _, v := range process.tasks {
		fmt.Println(v)
	}
}

// 根据流程类型创建流程实例
func CreateProcess(process_type int) *Process{
	var process = new(Process)
	process.process_type = process_type
	process.getNodes()
	process.current_node = process.getStartNode()
	process.createOrder()
	process.getTasks()
	return process
}

// 根据order_id获取流程实例
func GetProcess(order_id int) *Process{
	var process = new(Process)
	process.order = process.getOrderByID(int(order_id))
	process.process_type = process.order.process_type
	process.getNodes()
	process.current_node = process.getNodeByID(process.order.current_node)
	process.getTasks()
	return process
}

// 记录流程信息
func (process *Process) createOrder(){
	create_user := 100 
	create_date := "20140307145252"
	sql := "insert into `process_order` (process_type, current_node, create_user, create_time) values(?, ?, ?, ?);"

	db := lessgo.GetMySQL()
	defer db.Close()

	stmt, err := db.Prepare(sql)
	checkErr(err)
	res, err := stmt.Exec(process.process_type, process.current_node.id, create_user, create_date)    
	checkErr(err)
	order_id, err := res.LastInsertId() 
	process.order = process.getOrderByID(int(order_id))
	checkErr(err)

	process.CreateTask(222, "review pass", "pass")
}

// 记录节点动作
func (process *Process) CreateTask(actor int, action string, result string){
	sql := "insert `process_task` (order_id, node_id, actor, action, result) values(?, ?, ?, ?, ?)"

	db := lessgo.GetMySQL()
	defer db.Close()

	stmt, err := db.Prepare(sql)
	checkErr(err)
        res, err := stmt.Exec(process.order.id, process.current_node.id, actor, action, result) 
	checkErr(err)
	id, err := res.LastInsertId() 
	process.task_id = int(id)
        checkErr(err)

	var nextID int
        if result == "pass" { 
		nextID = process.current_node.pass_id 
	}else if result == "back"{
		nextID = process.current_node.pass_id 
	}
        process.updateCurrentNode(nextID)
}

// 更新当前节点
func (process *Process) updateCurrentNode(node_id int) {
	process.current_node = process.getNodeByID(node_id)

	db := lessgo.GetMySQL()
	defer db.Close()
	sql := "update `process_order` set current_node=? where id=?"
	stmt, err := db.Prepare(sql)
	checkErr(err)
        res, err := stmt.Exec(node_id, process.order.id) 
	checkErr(err)
	res.LastInsertId()
}

// 获取开始节点
func (process *Process) getStartNode() *node{
	db := lessgo.GetMySQL()
	defer db.Close()
	sql := "select * from `process_node` where process_type=? and node_type='start'"
	rows, err := db.Query(sql, process.process_type)
	checkErr(err)

	tmpNode := new(node)
	if rows.Next() {
		err := rows.Scan(&tmpNode.id, &tmpNode.process_type, &tmpNode.back_id, &tmpNode.pass_id, &tmpNode.node_type, &tmpNode.role) 
		checkErr(err)
	}
	return tmpNode
}

// 根据节点ID获取节点
func (process *Process) getNodeByID(node_id int) *node{
	db := lessgo.GetMySQL()
	defer db.Close()
	sql := "select * from `process_node` where id=?"
	rows, err := db.Query(sql, node_id)
	checkErr(err)

	tmpNode := new(node)
	if rows.Next() {
		err := rows.Scan(&tmpNode.id, &tmpNode.process_type, &tmpNode.back_id, &tmpNode.pass_id, &tmpNode.node_type, &tmpNode.role) 
		checkErr(err)
	}
	return tmpNode
}

// 根据流程类型获取所有节点
func (process *Process) getNodes() {
	db := lessgo.GetMySQL()
	defer db.Close()
	sql := "select * from `process_node` where process_type=?"
	rows, err := db.Query(sql, process.process_type)
	checkErr(err)
	for rows.Next() {
	         tmpNode := new(node)
		 err := rows.Scan(&tmpNode.id, &tmpNode.process_type, &tmpNode.back_id, &tmpNode.pass_id, &tmpNode.node_type, &tmpNode.role) 
		 checkErr(err)
		 process.nodes = append(process.nodes, tmpNode)
	}
}

// 根据order_id获取order
func (process *Process) getOrderByID(order_id int) *order{
	db := lessgo.GetMySQL()
	defer db.Close()
	sql := "select * from `process_order` where id=?"
	rows, err := db.Query(sql, order_id)
	checkErr(err)

	tmpOrder := new(order)
	if rows.Next() {
		err := rows.Scan(&tmpOrder.id, &tmpOrder.process_type, &tmpOrder.current_node, &tmpOrder.create_user, &tmpOrder.create_time) 
		checkErr(err)
	}
	return tmpOrder
}

// 获取动作列表
func (process *Process) getTasks(){
	db := lessgo.GetMySQL()
	defer db.Close()
	sql := "select * from `process_task` where order_id=?"
	rows, err := db.Query(sql, process.order.id)
	checkErr(err)
	for rows.Next() {
		tmpTask := new(task)
		err := rows.Scan(&tmpTask.id, &tmpTask.order_id, &tmpTask.node_id, &tmpTask.actor, &tmpTask.action, &tmpTask.result) 
		checkErr(err)
		process.tasks = append(process.tasks, tmpTask)
	}
}

func checkErr(err error) { 
	m := make(map[string]interface{}) 
	if err != nil {
		m["success"] = false
		m["code"] = 100 
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
//		commonlib.OutputJson(w, m, " ")
		return
	}
}
