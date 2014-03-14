define(function (require, exports, module) {
    var $ = require('jquery');
    require('juicer');
    require('timepicker');
    
    var currentEmployee    = {};
    var accountingSubjects = {};
    var headquartersId     = 9;
    var process = {
		init : function(){
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
