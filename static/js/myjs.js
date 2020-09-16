
$(function () {

    // $("input[name='startTime']").datetimepicker({
    //     minView: "hour", //  选择时间时，最小可以选择到那层；默认是‘hour’也可用0表示
    //     language: 'zh-CN', // 语言
    //     autoclose: true, //  true:选择时间后窗口自动关闭
    //     format: 'yyyy-mm-dd hh:ii:00', // 文本框时间格式，设置为0,最后时间格式为2017-03-23 17:15:00
    //     // todayBtn: true, // 如果此值为true 或 "linked"，则在日期时间选择器组件的底部显示一个 "Today" 按钮用以选择当前日期。
    //     startDate: new Date(),  // 窗口可选时间从今天开始
    //     // endDate: new Date()   // 窗口最大时间直至今天
    // })
    laydate.render({
        elem: '#laydate-choose-date',
        format: 'yyyy/MM/dd',
        theme: '#393D49',
        min: 0, //今天
        showBottom: false
    });
    laydate.render({
        elem: '#laydate-choose-time',
        type: 'time',
        btns: ['confirm'],
        theme: '#393D49',
    })

});