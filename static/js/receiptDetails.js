define(function (require, exports, module) {
    var $ = require('jquery');
    require('juicer');
    require('timepicker');

    // 基本信息对象
    var $divReceiptDetails = $('#divReceiptDetails');

    $divReceiptDetails.html('待添加数据列表');

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
	$('.editItem').bind('click', editItem);
	$('.saveItem').bind('click', saveItem);
	$('.deleteItem').bind('click', deleteItem);
    });

    function editItem()
    {
	$(this).prop('hidden', true);
	$(this).next().prop('hidden', false);
	$(this).parent().prevAll().each(function(){ $(this).find('[disabled]').prop('disabled', false); });

    }
    function saveItem()
    {
	$(this).prop('hidden', true);
	$(this).prev().prop('hidden', false)
	$(this).parent().prevAll().each(function(){
		$(this).find('select').prop('disabled', 'disabled');
		$(this).find('input').prop('disabled', 'disabled');
		$(this).find('textarea').prop('disabled', 'disabled');
	});
    }
    function deleteItem()
    {
	$(this).parent().parent().remove();
    }


});
