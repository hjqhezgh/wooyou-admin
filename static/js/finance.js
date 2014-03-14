define(function (require, exports, module) {
    var $ = require('jquery');
    require('juicer');
    require('bootstrap');
    // 基本信息对象
    var $divSubmissionAction = $('#divSubmissionAction');
    var $divSearch = $('#divSearch');
    var $divCenter = $('#divCenter');
    var $divHandleApply = $('#divHandleApply');
    var $divApplyDetails = $('#divApplyDetails');
    var $btnReimbursed = $('#btnReimbursed');
    var $btnPayment = $('#btnPayment');
    // 模板
    var handleApplyTemp = '<div style="display: inline;margin-right:5%;"><a href="/pendingReceiptDetails"><span class="text-error">${textPendingReceiptAmount}</span></a></div>' +
        '<div style="display: inline;"><a href="/completedReceiptDetails">${textCompletedReceiptAmount}</a></div>';
    var $WEEK_INFO = '<div><a href="/classifiedPendingReceiptList?status=P&week=${week}&code=${accountingSubjectCode}" title="未审核或异常"><span class="text-error">$${textPendingReceiptAmount}</span></a></div>' +
        '<div><a href="/classifiedCompletedReceiptList?status=A&week=${week}&code=${accountingSubjectCode}" title="财务和主管都通过"><span class="text-info">$${textCompletedReceiptAmount}</span></a></div>';
    var $TABLE = '<table class="table table-box" border=1>' +
        '<thead>' +
        '<tr class="text-center">' +
        '<th class="financeType">财务类型</th>' +
        '<th class="week1st">第一周</th>' +
        '<th class="week2nd">第二周</th>' +
        '<th class="week3rd">第三周</th>' +
        '<th class="week4th">第四周</th>' +
        '<th class="week5th">第五周</th>' +
        '<th class="week6th">第六周</th>' +
        '<th class="sub_total">总计</th>' +
        '</tr>' +
        '</thead>' +
        '<tbody>$${tBody}</tbody>' +
        '</table>';
    var $TR = '<tr>' +
        '<td class="financeType">$${financeType}</td>' +
        '<td class="week1st monetary">$${week1st}</td>' +
        '<td class="week2nd monetary">$${week2nd}</td>' +
        '<td class="week3rd monetary">$${week3rd}</td>' +
        '<td class="week4th monetary">$${week4th}</td>' +
        '<td class="week5th monetary">$${week5th}</td>' +
        '<td class="week6th monetary">$${week6th}</td>' +
        '<td class="sub_total monetary">$${subTotal}</td>' +
        '</tr>';

    $.get("/getRoleCodes.json", {}, function (data) {
        if (data.success) {
            var roleCodes = data.datas;
            render(roleCodes);
//            bindEvent();
        } else {
            alert(data.msg);
        }
    }, 'json');

    function render(roleCodes) {
        if (isExistInArray(roleCodes, 'cd')) {
            $divSubmissionAction.show();
            $divCenter.css('display', 'none');
        }
        if (isExistInArray(roleCodes, 'yyzj')) {
            $divSubmissionAction.hide();
            $divCenter.css('display', 'inline-block');
        }
        $divCenter.append('<select><option>A中心</option><option>B中心</option></select>');

        // 处理信息汇总
        $.get("/getHandleApplyInfo.json", {}, function (data) {
            if (data.success) {
                var handleApplyInfo = data.datas;
                var pendingReceiptAmount = handleApplyInfo.pendingReceiptAmount;
                var completedReceiptAmount = handleApplyInfo.completedReceiptAmount;
                // 获取数据
                $divHandleApply.append(juicer(handleApplyTemp, {
                    textPendingReceiptAmount: pendingReceiptAmount + "(未处理单据)",
                    textCompletedReceiptAmount: completedReceiptAmount + "(已处理单据)"
                }));
            } else {
                alert('获取数据出错！');
            }
        }, 'json');
        // 详细信息列表
        $.get("/getReceiptDetails.json", {}, function (data) {
            if (data.success) {
                var applyDetailsMap = getApplyDetailsMap(data.datas);
                var $TRs = getApplyDetailsDisplayList(applyDetailsMap);
                renderApplyDetailsList($TRs);
            } else {
                alert('获取数据出错！');
            }
        }, 'json')
    }

    function isExistInArray(arrData, code) {
        var isExist = false;
        for (var i in arrData) {
            if (code.toLowerCase() == arrData[i].toLowerCase()) {
                isExist = true;
                break;
            }
        }

        return isExist;
    }

    function renderApplyDetailsList($TRs) {
        $divApplyDetails.append(juicer($TABLE, {
            tBody: $TRs
        }));
        var weeksNumOfMonth = getWeeksNumOfMonth(new Date());
        var $WEEK5TH = $('.week5th');
        var $WEEK6TH = $('.week6th');
        if (weeksNumOfMonth == 5) {
            $WEEK6TH.hide();
        } else if (weeksNumOfMonth == 4) {
            $WEEK5TH.hide();
            $WEEK6TH.hide();
        }
        var $TH_CELL = $('table > thead > tr > th');
        var $TD_MONETARY = $('table .monetary');
        $TH_CELL.css('text-align', 'center');
        $TD_MONETARY.css('text-align', 'right');
        $('table, tr, td, th').css('border', '1px solid gray');
    }

    function getApplyDetailsMap(receiptDetails) {
        var applyDetailsMap = {};
        for (var i = 0; i < receiptDetails.length; i++) {
            var currentReceiptDetail = receiptDetails[i];
//                    var receiptId = currentReceiptDetail.id;
            var accountingSubjectCode = currentReceiptDetail.accountingSubjectCode;
            var accountingSubjectName = currentReceiptDetail.accountingSubjectName;
            var receiptAmount = currentReceiptDetail.receiptAmount;
            var applyDate = currentReceiptDetail.applyDate;
            var approveStatus = currentReceiptDetail.approveStatus;

            var weekOfMonth = getWeekOfMonth(new Date(applyDate), false).toString();

            var applyDetailKey = accountingSubjectCode;
            var applyDetailValue = applyDetailsMap[applyDetailKey];
            if (applyDetailValue == null) {
                var applyDetail = new Object();
                applyDetail.accountingSubjectCode = accountingSubjectCode;
                applyDetail.accountingSubjectName = accountingSubjectName;
                applyDetail.week1stPendingReceiptAmount = 0;
                applyDetail.week1stCompletedReceiptAmount = 0;
                applyDetail.week2ndPendingReceiptAmount = 0;
                applyDetail.week2ndCompletedReceiptAmount = 0;
                applyDetail.week3rdPendingReceiptAmount = 0;
                applyDetail.week3rdCompletedReceiptAmount = 0;
                applyDetail.week4thPendingReceiptAmount = 0;
                applyDetail.week4thCompletedReceiptAmount = 0;
                applyDetail.week5thPendingReceiptAmount = 0;
                applyDetail.week5thCompletedReceiptAmount = 0;
                applyDetail.week6thPendingReceiptAmount = 0;
                applyDetail.week6thCompletedReceiptAmount = 0;
                applyDetail.subTotalPendingReceiptAmount = 0;
                applyDetail.subTotalCompletedReceiptAmount = 0;
                if (approveStatus == "A") {
                    switch (weekOfMonth) {
                        case "1":
                            applyDetail.week1stCompletedReceiptAmount = parseFloat(receiptAmount);
                            break;
                        case "2":
                            applyDetail.week2ndCompletedReceiptAmount = parseFloat(receiptAmount);
                            break;
                        case "3":
                            applyDetail.week3rdCompletedReceiptAmount = parseFloat(receiptAmount);
                            break;
                        case "4":
                            applyDetail.week4thCompletedReceiptAmount = parseFloat(receiptAmount);
                            break;
                        case "5":
                            applyDetail.week5thCompletedReceiptAmount = parseFloat(receiptAmount);
                            break;
                        case "6":
                            applyDetail.week6thCompletedReceiptAmount = parseFloat(receiptAmount);
                            break;
                        default:
                            break;
                    }
                    applyDetail.subTotalCompletedReceiptAmount = parseFloat(receiptAmount);
                } else {
                    switch (weekOfMonth) {
                        case "1":
                            applyDetail.week1stPendingReceiptAmount = parseFloat(receiptAmount);
                            break;
                        case "2":
                            applyDetail.week2ndPendingReceiptAmount = parseFloat(receiptAmount);
                            break;
                        case "3":
                            applyDetail.week3rdPendingReceiptAmount = parseFloat(receiptAmount);
                            break;
                        case "4":
                            applyDetail.week4thPendingReceiptAmount = parseFloat(receiptAmount);
                            break;
                        case "5":
                            applyDetail.week5thPendingReceiptAmount = parseFloat(receiptAmount);
                            break;
                        case "6":
                            applyDetail.week6thPendingReceiptAmount = parseFloat(receiptAmount);
                            break;
                        default:
                            break;
                    }
                    applyDetail.subTotalPendingReceiptAmount = parseFloat(receiptAmount);
                }

                applyDetailsMap[applyDetailKey] = applyDetail;
            } else {
                if (approveStatus == "A") {
                    switch (weekOfMonth) {
                        case "1":
                            applyDetailsMap[applyDetailKey].week1stCompletedReceiptAmount += parseFloat(receiptAmount);
                            break;
                        case "2":
                            applyDetailsMap[applyDetailKey].week2ndCompletedReceiptAmount += parseFloat(receiptAmount);
                            break;
                        case "3":
                            applyDetailsMap[applyDetailKey].week3rdCompletedReceiptAmount += parseFloat(receiptAmount);
                            break;
                        case "4":
                            applyDetailsMap[applyDetailKey].week4thCompletedReceiptAmount += parseFloat(receiptAmount);
                            break;
                        case "5":
                            applyDetailsMap[applyDetailKey].week5thCompletedReceiptAmount += parseFloat(receiptAmount);
                            break;
                        case "6":
                            applyDetailsMap[applyDetailKey].week6thCompletedReceiptAmount += parseFloat(receiptAmount);
                            break;
                        default:
                            break;
                    }
                    applyDetailsMap[applyDetailKey].subTotalCompletedReceiptAmount += parseFloat(receiptAmount);
                } else {
                    switch (weekOfMonth) {
                        case "1":
                            applyDetailsMap[applyDetailKey].week1stPendingReceiptAmount += parseFloat(receiptAmount);
                            break;
                        case "2":
                            applyDetailsMap[applyDetailKey].week2ndPendingReceiptAmount += parseFloat(receiptAmount);
                            break;
                        case "3":
                            applyDetailsMap[applyDetailKey].week3rdPendingReceiptAmount += parseFloat(receiptAmount);
                            break;
                        case "4":
                            applyDetailsMap[applyDetailKey].week4thPendingReceiptAmount += parseFloat(receiptAmount);
                            break;
                        case "5":
                            applyDetailsMap[applyDetailKey].week5thPendingReceiptAmount += parseFloat(receiptAmount);
                            break;
                        case "6":
                            applyDetailsMap[applyDetailKey].week6thPendingReceiptAmount += parseFloat(receiptAmount);
                            break;
                        default:
                            break;
                    }
                    applyDetailsMap[applyDetailKey].subTotalPendingReceiptAmount += parseFloat(receiptAmount);
                }
            }
        }

        return applyDetailsMap;
    }

    function getApplyDetailsDisplayList(applyDetailsMap) {
        var $TRs = '';
        var weekTotal = new Object();
        weekTotal.accountingSubjectCode = "WT";
        weekTotal.accountingSubjectName = "总计";
        weekTotal.week1stPendingReceiptAmount = 0;
        weekTotal.week1stCompletedReceiptAmount = 0;
        weekTotal.week2ndPendingReceiptAmount = 0;
        weekTotal.week2ndCompletedReceiptAmount = 0;
        weekTotal.week3rdPendingReceiptAmount = 0;
        weekTotal.week3rdCompletedReceiptAmount = 0;
        weekTotal.week4thPendingReceiptAmount = 0;
        weekTotal.week4thCompletedReceiptAmount = 0;
        weekTotal.week5thPendingReceiptAmount = 0;
        weekTotal.week5thCompletedReceiptAmount = 0;
        weekTotal.week6thPendingReceiptAmount = 0;
        weekTotal.week6thCompletedReceiptAmount = 0;
        weekTotal.subTotalPendingReceiptAmount = 0;
        weekTotal.subTotalCompletedReceiptAmount = 0;
        for (var key in applyDetailsMap) {
            // 列数据统计
            weekTotal.week1stPendingReceiptAmount += applyDetailsMap[key].week1stPendingReceiptAmount;
            weekTotal.week1stCompletedReceiptAmount += applyDetailsMap[key].week1stCompletedReceiptAmount;
            weekTotal.week2ndPendingReceiptAmount += applyDetailsMap[key].week2ndPendingReceiptAmount;
            weekTotal.week2ndCompletedReceiptAmount += applyDetailsMap[key].week2ndCompletedReceiptAmount;
            weekTotal.week3rdPendingReceiptAmount += applyDetailsMap[key].week3rdPendingReceiptAmount;
            weekTotal.week3rdCompletedReceiptAmount += applyDetailsMap[key].week3rdCompletedReceiptAmount;
            weekTotal.week4thPendingReceiptAmount += applyDetailsMap[key].week4thPendingReceiptAmount;
            weekTotal.week4thCompletedReceiptAmount += applyDetailsMap[key].week4thCompletedReceiptAmount;
            weekTotal.week5thPendingReceiptAmount += applyDetailsMap[key].week5thPendingReceiptAmount;
            weekTotal.week5thCompletedReceiptAmount += applyDetailsMap[key].week5thCompletedReceiptAmount;
            weekTotal.week6thPendingReceiptAmount += applyDetailsMap[key].week6thPendingReceiptAmount;
            weekTotal.week6thCompletedReceiptAmount += applyDetailsMap[key].week6thCompletedReceiptAmount;
            weekTotal.subTotalPendingReceiptAmount += applyDetailsMap[key].subTotalPendingReceiptAmount;
            weekTotal.subTotalCompletedReceiptAmount += applyDetailsMap[key].subTotalCompletedReceiptAmount;
            var financeType = applyDetailsMap[key].accountingSubjectName;
            var week1st = juicer($WEEK_INFO, {
                week: "1",
                accountingSubjectCode: applyDetailsMap[key].accountingSubjectCode,
                textPendingReceiptAmount: formatCurrency(applyDetailsMap[key].week1stPendingReceiptAmount),
                textCompletedReceiptAmount: formatCurrency(applyDetailsMap[key].week1stCompletedReceiptAmount)
            });
            var week2nd = juicer($WEEK_INFO, {
                week: "2",
                accountingSubjectCode: applyDetailsMap[key].accountingSubjectCode,
                textPendingReceiptAmount: formatCurrency(applyDetailsMap[key].week2ndPendingReceiptAmount),
                textCompletedReceiptAmount: formatCurrency(applyDetailsMap[key].week2ndCompletedReceiptAmount)
            });
            var week3rd = juicer($WEEK_INFO, {
                week: "3",
                accountingSubjectCode: applyDetailsMap[key].accountingSubjectCode,
                textPendingReceiptAmount: formatCurrency(applyDetailsMap[key].week3rdPendingReceiptAmount),
                textCompletedReceiptAmount: formatCurrency(applyDetailsMap[key].week3rdCompletedReceiptAmount)
            });
            var week4th = juicer($WEEK_INFO, {
                week: "4",
                accountingSubjectCode: applyDetailsMap[key].accountingSubjectCode,
                textPendingReceiptAmount: formatCurrency(applyDetailsMap[key].week4thPendingReceiptAmount),
                textCompletedReceiptAmount: formatCurrency(applyDetailsMap[key].week4thCompletedReceiptAmount)
            });
            var week5th = juicer($WEEK_INFO, {
                week: "5",
                accountingSubjectCode: applyDetailsMap[key].accountingSubjectCode,
                textPendingReceiptAmount: formatCurrency(applyDetailsMap[key].week5thPendingReceiptAmount),
                textCompletedReceiptAmount: formatCurrency(applyDetailsMap[key].week5thCompletedReceiptAmount)
            });
            var week6th = juicer($WEEK_INFO, {
                week: "6",
                accountingSubjectCode: applyDetailsMap[key].accountingSubjectCode,
                textPendingReceiptAmount: formatCurrency(applyDetailsMap[key].week6thPendingReceiptAmount),
                textCompletedReceiptAmount: formatCurrency(applyDetailsMap[key].week6thCompletedReceiptAmount)
            });
            var subTotal = juicer($WEEK_INFO, {
                week: "0",
                accountingSubjectCode: applyDetailsMap[key].accountingSubjectCode,
                textPendingReceiptAmount: formatCurrency(applyDetailsMap[key].subTotalPendingReceiptAmount),
                textCompletedReceiptAmount: formatCurrency(applyDetailsMap[key].subTotalCompletedReceiptAmount)
            });

            var applyDetailRow = juicer($TR, {
                financeType: financeType,
                week1st: week1st,
                week2nd: week2nd,
                week3rd: week3rd,
                week4th: week4th,
                week5th: week5th,
                week6th: week6th,
                subTotal: subTotal
            })

            $TRs += applyDetailRow;
        }

        $TRs += renderWeekTotalRow(weekTotal);

        return $TRs;
    }
    
    function renderWeekTotalRow(weekTotal) {
        var week1st = juicer($WEEK_INFO, {
            week: "1",
            accountingSubjectCode: weekTotal.accountingSubjectCode,
            textPendingReceiptAmount: formatCurrency(weekTotal.week1stPendingReceiptAmount),
            textCompletedReceiptAmount: formatCurrency(weekTotal.week1stCompletedReceiptAmount)
        });
        var week2nd = juicer($WEEK_INFO, {
            week: "2",
            accountingSubjectCode: weekTotal.accountingSubjectCode,
            textPendingReceiptAmount: formatCurrency(weekTotal.week2ndPendingReceiptAmount),
            textCompletedReceiptAmount: formatCurrency(weekTotal.week2ndCompletedReceiptAmount)
        });
        var week3rd = juicer($WEEK_INFO, {
            week: "3",
            accountingSubjectCode: weekTotal.accountingSubjectCode,
            textPendingReceiptAmount: formatCurrency(weekTotal.week3rdPendingReceiptAmount),
            textCompletedReceiptAmount: formatCurrency(weekTotal.week3rdCompletedReceiptAmount)
        });
        var week4th = juicer($WEEK_INFO, {
            week: "4",
            accountingSubjectCode: weekTotal.accountingSubjectCode,
            textPendingReceiptAmount: formatCurrency(weekTotal.week4thPendingReceiptAmount),
            textCompletedReceiptAmount: formatCurrency(weekTotal.week4thCompletedReceiptAmount)
        });
        var week5th = juicer($WEEK_INFO, {
            week: "5",
            accountingSubjectCode: weekTotal.accountingSubjectCode,
            textPendingReceiptAmount: formatCurrency(weekTotal.week5thPendingReceiptAmount),
            textCompletedReceiptAmount: formatCurrency(weekTotal.week5thCompletedReceiptAmount)
        });
        var week6th = juicer($WEEK_INFO, {
            week: "6",
            accountingSubjectCode: weekTotal.accountingSubjectCode,
            textPendingReceiptAmount: formatCurrency(weekTotal.week6thPendingReceiptAmount),
            textCompletedReceiptAmount: formatCurrency(weekTotal.week6thCompletedReceiptAmount)
        });
        var subTotal = juicer($WEEK_INFO, {
            week: "0",
            accountingSubjectCode: weekTotal.accountingSubjectCode,
            textPendingReceiptAmount: formatCurrency(weekTotal.subTotalPendingReceiptAmount),
            textCompletedReceiptAmount: formatCurrency(weekTotal.subTotalCompletedReceiptAmount)
        });

        var applyDetailRow = juicer($TR, {
            financeType: weekTotal.accountingSubjectName,
            week1st: week1st,
            week2nd: week2nd,
            week3rd: week3rd,
            week4th: week4th,
            week5th: week5th,
            week6th: week6th,
            subTotal: subTotal
        })

        return applyDetailRow;
    }
    /**
     * 获取指定日期在当前月份是第几周
     * date：指定日期
     * isBeginSunday: 周从星期日开始
      */
    function getWeekOfMonth(date, isBeginSunday) {
        var firstDateOfMonth = new Date(date.getFullYear(), date.getMonth(), 1);
        var firstDateOfWeek = firstDateOfMonth.getDay();
        var dayOfMonth = date.getDate();
        var weekOfMonth = 0;
        if (isBeginSunday) {
            weekOfMonth = parseInt((dayOfMonth + firstDateOfWeek - 1) / 7) + 1;
        } else {
            // 1号正好是周日
            if (firstDateOfWeek == 0) {
                firstDateOfWeek = 7;
            }
            weekOfMonth = parseInt((dayOfMonth + firstDateOfWeek - 2) / 7) + 1;
        }

        return weekOfMonth;
    }
    // 获取指定时间的月份的天数
    function getDaysNumOfMonth(date) {
        return (new Date(date.getFullYear(), date.getMonth() + 1, 0)).getDate();
    }
    // 获取指定时间的月份的周数
    function getWeeksNumOfMonth(date) {
        return getWeekOfMonth(new Date(date.getFullYear(), date.getMonth(), getDaysNumOfMonth(date)));
    }
    // 转换数字格式为 123,456,789.01 格式
    function formatCurrency(num) {
        num = num.toString().replace(/\$|\,/g, '');
        if (isNaN(num)) {
            num = "0";
        }
        var sign = (num == (num = Math.abs(num)));
        num = Math.floor(num * 100 + 0.50000000001);
        var cents = num % 100;
        num = Math.floor(num / 100).toString();

        if (cents < 10) {
            cents = "0" + cents;
        }

        for (var i = 0; i < Math.floor((num.length - (1 + i)) / 3); i++) {
            num = num.substring(0, num.length - (4 * i + 3)) + ',' + num.substring(num.length - (4 * i + 3));
        }

        if (cents == "00") {
            return ((sign) ? '' : '-') + num
        } else {
            return (((sign) ? '' : '-') + num + '.' + cents);
        }
    }

    $(document).on('click','a[data-action=openIframeWindow]',function(e){
        e.preventDefault();

        var url = $(this).attr('href');

        $.openIframeWindow({
            url : url,
            parentComponent : "",
            parentWindowName : window.name
        },e);
    });

});
