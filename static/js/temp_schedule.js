define(function (require, exports, module) {
    var $ = require('jquery');
    require('/js/lib/superTables/superTables.js');
    require('/js/lib/superTables/superTables.css');
    require('timepicker');
    require('lessgo-window');
    require('juicer');

    var dataUrl = '/web/class_schedule_detail/temp_data.json';

    var container = $('.fakeContainer');

    var process = {
        init : function(){

            $('select[name=centerId]').html('<option value="">全部</option></select>');

            $.get('/web/center/alldata',function(data){
                var optionTemp = '<option value="${value}">${desc}</option>';

                if(data.success){
                    var objects = data.datas;

                    for(var i=0;i<objects.length;i++){
                        $('select[name=centerId]').append(juicer(optionTemp,{
                            value:objects[i]["cid"],
                            desc:objects[i]["name"]
                        }));
                    }

                } else{
                    alert(data.msg);
                }
            },'json');

            process.render();

            $('input[name=start_time]').zIndex(10);

            $('input[name=start_time]').datepicker({
                dateFormat : 'yy-mm-dd'
            });

            process.bind();
        },
        render : function(type){

            container.html('<table border="1" id="schedule" width="100%" border="0" cellpadding="0" cellspacing="0"></table>');
            var allDataTable = $('#schedule');
            $.get(dataUrl,{
                "date" :  $('[name=firstDayOfWeek]').val(),
                "centerId" :  $('[name=centerId]').val(),
                "type" : type
            },function(data){
                if(data.success){
                    $('[name=firstDayOfWeek]').val(data.firstDayOfWeek);

                    var weekdates = data.weekdates;
                    var roomLength = data.rooms.length;
                    allDataTable.append('<tr></tr>');
                    allDataTable.find('tr').append('<th rowspan="2">&nbsp;</th>');
                    for(var i=0;i<weekdates.length;i++){
                        allDataTable.find('tr').append('<th colspan="'+roomLength+'">'+weekdates[i].week+'('+weekdates[i].date+')</th>');
                    }

                    allDataTable.append('<tr></tr>');
                    for(var i=0;i<7;i++){
                        for(var j=0;j<data.rooms.length;j++){
                            allDataTable.find('tr:last').append('<td><div class="room">'+data.rooms[j].name+'</div></td>');
                        }
                    }

                    for(var i=0;i<data.times.length;i++){
                        allDataTable.append('<tr></tr>');
                        allDataTable.find('tr:last').append('<td>'+data.times[i].name+'</td>');

                        var num = roomLength * 7;
                        for(var j=0;j<num;j++){
                            var schedule = process.findSchedule(data.schedules,data.times[i].id,data.rooms[j%roomLength].id,Math.ceil((j+0.01)/roomLength));
                            process.renderCell(schedule,data.times[i].id,data.rooms[j%roomLength].id,Math.ceil((j+0.01)/roomLength),weekdates[Math.ceil((j+0.01)/roomLength)-1].date);
                        }

                    }

                    new superTable("schedule", {
                        cssSkin : "sDefault",
                        fixedCols : 1,
                        headerRows : 2
                    });

                }else{
                    alert(data.msg)
                }

            },'json');

        },
        findSchedule: function(schedules,timeId,roomId,week){
            for(var k=0;k<schedules.length;k++){
                var schedule = schedules[k];
                if(schedule.timeId==timeId && schedule.roomId==roomId && schedule.week == week){
                    return schedule;
                }
            }

            return null;
        },
        renderCell : function(schedule,timeId,roomId,week,date){
            var allDataTable = $('#schedule');
            if(schedule){
                if(schedule.isNormal == 1){
                    allDataTable.find('tr:last').append('<td><div class="schedule-detail normal"></div></td>');
                    allDataTable.find('div:last').append('<p>'+schedule.name+'</p>');
                    allDataTable.find('div:last').append('<p>老师：'+schedule.teacher+'</p>');
                    allDataTable.find('div:last').append('<p>助教：'+schedule.assistant+'</p>');
                    allDataTable.find('div:last').append('<p>人数：'+schedule.signNum+"/"+schedule.personNum+'</p>');
                    allDataTable.find('div:last').append('<p><a href="/web/wyclass/manageChildForNormal?centerId-eq='+schedule.centerId+'&scheduleId='+schedule.id+'" data-action="openIframeWindow">分配学生</a></p>');
                    allDataTable.find('div:last').append('<p><a href="#" data-value="'+schedule.id+'" data-action="deleteTmpSchedule">删</a>&nbsp;<a href="/web/class_schedule_detail/modify_normal?scheduleId='+schedule.id+'" data-value="'+schedule.id+'" data-action="openIframeWindow">改</a></p>');
                }else if(schedule.isNormal == 2){
                    allDataTable.find('tr:last').append('<td><div class="schedule-detail foronce"></div></td>');
                    if(schedule.code){
                        allDataTable.find('div:last').append('<p style="color:red;">'+schedule.name+'('+schedule.code+')</p>');
                    }else{
                        allDataTable.find('div:last').append('<p style="color:red;">'+schedule.name+'</p>');
                    }
                    allDataTable.find('div:last').append('<p>当前人数：'+schedule.currentTMKPersonNum+"/"+schedule.personNum+'</p>');
                    allDataTable.find('div:last').append('<p>签到人数：'+schedule.signNum+'</p>');
                    allDataTable.find('div:last').append('<p><a href="/web/wyclass/manageChild?classId='+schedule.classId+'&centerId-eq='+schedule.centerId+'&scheduleId='+schedule.id+'" data-action="openIframeWindow">学生管理</a></p>');
                    allDataTable.find('div:last').append('<p></p>');
                }
            }else{
                var emptyTdTmp = '<td><div class="schedule-detail"><p>无</p></div></td>'

                allDataTable.find('tr:last').append(juicer(emptyTdTmp,{
                    roomId: roomId,
                    timeId:timeId,
                    week:week,
                    date:date
                }));

            }
        },
        bind: function () {
            var allDataTable = $('#schedule');

            $('a[data-action=preWeek]').click(function(e){
                e.preventDefault();
                process.render("pre");
            });

            $('a[data-action=nextWeek]').click(function(e){
                e.preventDefault();
                process.render("next");
            });

            $('a[data-action=search]').click(function(e){
                e.preventDefault();
                if($('input[name=start_time]').val()){
                    $('[name=firstDayOfWeek]').val($('input[name=start_time]').val().replace(new RegExp("-","g"), "")+"000000");
                    process.render();
                }
            });

            $('a[data-action=reset]').click(function(e){
                e.preventDefault();
                $('input[name=start_time]').val('')
            });

            $('a[data-action=createWeekSchedule]').click(function(e){
                e.preventDefault();

                $.post('/web/class_schedule_detail/createWeekSchedule.json',{
                    firstDayOfWeek : $('[name=firstDayOfWeek]').val()
                },function(data){
                    if(data.success){
                        process.render();
                    }else{
                        alert(data.msg);
                    }
                },'json');
            });

            allDataTable.on('click','a[data-action=openIframeWindow]',function(e){
                e.preventDefault();

                var url = $(this).attr('href');

                $.openIframeWindow({
                    url : url,
                    parentComponent : "",
                    parentWindowName : window.name
                },e);
            });

            $(document).on('click','a[data-action=deleteTmpSchedule]',function(e){
                e.preventDefault();

                var sid = $(this).attr('data-value');

                if(confirm("确认删除此课表吗?")){
                    $.post('/web/schedule_detail/deleteSingle.json',{scheduleId:sid},function(data){
                        if(data.success){
                            process.render();
                        }else{
                            alert(data.msg);
                        }
                    },'json');
                }
            });
        }
    }

    process.init();

    window.process = process;
});