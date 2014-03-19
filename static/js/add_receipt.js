define(function (require, exports, module) {
    var $ = require('jquery');
    require('juicer');
    require('timepicker');
    require('/js/lib/ztree/jquery.ztree.all-3.5.min.js');
    require('/js/lib/ztree/zTreeStyle.css');
	var setting = { view: { dblClickExpand: false }, data: { simpleData: { enable: true } }, callback: { beforeClick: beforeClick, onClick: onClick } };

		var zNodes =[
			{id:1, pId:0, name:"北京", accounting_subject_code:"4567"},
			{id:2, pId:0, name:"天津"},
			{id:3, pId:0, name:"上海"},
			{id:6, pId:0, name:"重庆"},
			{id:4400, pId:0, name:"河北省", open:true},
			{id:440001, pId:4400, name:"石家庄"},
			{id:440002, pId:4400, name:"保定"},
			{id:440003, pId:4400, name:"邯郸"},
			{id:440004, pId:4400, name:"承德"},
			{id:5, pId:0, name:"广东省", open:true},
			{id:51, pId:5, name:"广州"},
			{id:52, pId:5, name:"深圳"},
			{id:53, pId:5, name:"东莞"},
			{id:54, pId:5, name:"佛山"},
			{id:6, pId:0, name:"福建省", open:true},
			{id:61, pId:6, name:"福州"},
			{id:62, pId:6, name:"厦门"},
			{id:63, pId:6, name:"泉州"},
			{id:64, pId:6, name:"三明"}
		 ];

		function beforeClick(treeId, treeNode) {
			var check = (treeNode && !treeNode.isParent);
			if (!check) alert("只能选择子节点...");
			return check;
		}
		
		function onClick(e, treeId, treeNode) {
			var zTree = $.fn.zTree.getZTreeObj("treeDemo"),
			nodes = zTree.getSelectedNodes(),
			v = "";
			code = "";
			nodes.sort(function compare(a,b){return a.id-b.id;});
			for (var i=0, l=nodes.length; i<l; i++) {
				v += nodes[i].name + ",";
				code = nodes[i].code;
			}
			if (v.length > 0 ) v = v.substring(0, v.length-1);
			var cityObj = $(".ztreeInput");
			cityObj.attr("value", v);
//			cityObj.attr("name", accounting_subject_code);
			hideMenu();
		}

		function showMenu() {
			var cityObj = $(".ztreeInput");
			var cityOffset = $(".ztreeInput").offset();
			$("#menuContent").css({left:cityOffset.left + "px", top:cityOffset.top + cityObj.outerHeight() + "px"}).slideDown("fast");

			$("body").bind("mousedown", onBodyDown);
		}
		function hideMenu() {
			$("#menuContent").fadeOut("fast");
			$("body").unbind("mousedown", onBodyDown);
		}
		function onBodyDown(event) {
			if (!(event.target.id == "ztreeInput" || event.target.id == "menuContent" || $(event.target).parents("#menuContent").length>0)) {
				hideMenu();
			}
		}
    
    var currentEmployee    = {};
    var accountingSubjects = {};
    var headquartersId     = 9;
    var process = {
		init : function(){
			$.fn.zTree.init($("#treeDemo"), setting, zNodes);
			process.getCurrentEmployee();
			process.setDepartments();
			process.setAccountingSubject();
			process.bind();
		},
		getCurrentEmployee : function(){
			$.get("/getCurrentEmployee.json", function (data) {
				if(data.success){
					currentEmployee = data.datas;
				}else{
					currentEmployee.centerId = headquartersId; 
				}
					process.setCenters();
				}, 'json');
		},
		setCenters : function() {
			$.get('/web/center/alldata',function(data){
				var commondOption  = '<option value="${value}">${desc}</option>';
				var selectedOption = '<option value="${value}" selected>${desc}</option>';
				if(data.success){
					var objects = data.datas;
					for(var i=0;i<objects.length;i++){
						var currentOption = objects[i]["cid"] == currentEmployee.centerId ? selectedOption : commondOption;
						$('select[name=center_id]').append(juicer(currentOption,{ value:objects[i]["cid"], desc:objects[i]["name"] }));
					}
				}else{
					alert(data.msg);
				}

				process.setCenterRelated($('select[name=center_id]').val());
				$('select[name=center_id]').change(function(){
					centerId = parseInt($(this).val());
					process.setCenterRelated(centerId);
				});
			},'json');
		},
		setCenterRelated : function(centerId){
			if(centerId == headquartersId) $('select[name=department_id]').parent().parent().prop('hidden', true);
			if(centerId != headquartersId) $('select[name=department_id]').parent().parent().prop('hidden', false);
		},
		setDepartments : function(){
			$.get('/web/department/alldata',function(data){
				var optionTemp = '<option value="${value}">${desc}</option>';
				if(data.success){
					var objects = data.datas;
					for(var i=0;i<objects.length;i++){
						$('select[name=department_id]').append(juicer(optionTemp,{
							value:objects[i]["id"],
							desc:objects[i]["name"]
						}));
					}
				} else{
					alert(data.msg);
				}
			},'json');
		},
		setAccountingSubject : function(){
			$.get('/web/accounting_subject/alldata',function(data){
				var option  = '<option value="${value}">${desc}</option>';
				if(data.success){
					accountingSubjects = data.datas;
					for(var i=0; i < accountingSubjects.length; i++){
						if(accountingSubjects[i]["code"].length == 4){
							$('select#accounting_subject').append(juicer(option,{ value:accountingSubjects[i]["code"], desc:accountingSubjects[i]["name"] }));
						}
					}

				}else{
					alert(data.msg);
				}
			}, 'json');
		},
		editItem : function() {
			$(this).prop('hidden', true);
			$(this).next().prop('hidden', false);
			$(this).parent().prevAll().each(function(){ $(this).find('[disabled]').prop('disabled', false); });
		},
		saveItem : function(){
			$(this).prop('hidden', true);
			$(this).prev().prop('hidden', false)
			$(this).parent().prevAll().each(function(){
				$(this).find('select').prop('disabled', 'disabled');
				$(this).find('input').prop('disabled', 'disabled');
				$(this).find('textarea').prop('disabled', 'disabled');
			});
		},
		deleteItem : function(){
			$(this).parent().parent().remove();
		},
		bind : function(){
			$('input.date').datepicker({ dateFormat : 'yy-mm-dd' });
			$('input.ztreeInput').bind('click', showMenu); 
	

			var itemIndex = 1;
			$('#addItem').click(function(e){
				var itemTemplate = "<tr id='itemTemplate' hidden='hidden'>" + $('#itemTemplate').html() + "</tr>";

				$('#itemTemplate [name="use_type[]"]').val($('#newItem [name=data1]').val());
				$('#itemTemplate [name="certificate_quantity[]"]').val($('#newItem [name=data2]').val());
				$('#itemTemplate [name="amount[]"]').val($('#newItem [name=data3]').val());
				$('#itemTemplate [name="accounting_periods_start[]"]').val($('#newItem [name=data4]').val());
				$('#itemTemplate [name="accounting_periods_end[]"]').val($('#newItem [name=data5]').val());
				$('#itemTemplate [name="remark[]"]').val($('#newItem [name=data6]').val());

				$('#itemTemplate').removeAttr('hidden');
				$('#itemTemplate').attr('id', 'item' + itemIndex);
				$("#items tbody").prepend(itemTemplate);
				itemIndex++;

				$('#newItem input').val('');
				$('#newItem textarea').val('');
				$('.editItem').bind('click', process.editItem);
				$('.saveItem').bind('click', process.saveItem);
				$('.deleteItem').bind('click', process.deleteItem);
			});
		}
    }
    process.init();
    window.process = process;
});
