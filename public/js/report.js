var goGetChart;

$(function() {
    loadOverview();
    goGetChart = new Chart(document.getElementById("goGetChart").getContext("2d"), {
        type: "line",
        options: {
            scales: {
                yAxes: [{
                    stacked: true
                }]
            }
        }
    });
    loadInfo();
    loadDomains();

    $('#reportFilterForm select[name="domain_id"]').on('change', function() {
        loadPackages($(this).val());
    });

    $('#reportFilterForm').submit(function() {
        loadOverview();
        loadInfo();
        return false;
    });
});

function loadDomains() {
    $.get("/api/domains", {}, function(resp) {
        var ele = $('#reportFilterForm select[name="domain_id"]')
        var options = ['<option value="">Domain</option>'];
        for (var i = 0; i < resp.data.length; i ++) {
            options.push('<option value="' + resp.data[i].id + '">' + resp.data[i].name + '</option>')
        }
        ele.empty().html(options.join(""))
    }, 'json');
}

function loadPackages(domainID) {
    $.get("/api/packages", {
        domain_id: domainID,
    }, function(resp) {
        var ele = $('#reportFilterForm select[name="package_id"]')
        var options = ['<option value="">Package</option>'];
        for (var i = 0; i < resp.data.length; i ++) {
            options.push('<option value="' + resp.data[i].id + '">' + resp.data[i].path + '</option>')
        }
        ele.empty().html(options.join(""))
    }, 'json');
}

function loadOverview() {
    var loading = '<i class="fas fa-fw fa-circle-notch fa-spin"></i>';
    $('#overviewToday').html(loading);
    $('#overviewYesterday').html(loading);
    $('#overviewLastSevenDays').html(loading);
    $('#overviewLastThirtyDays').html(loading);
    $('#overviewTotal').html(loading);
    $.get("/report/overview", $('#reportFilterForm').serialize(), function(resp) {
        $('#overviewToday').text(resp.data.today);
        $('#overviewYesterday').text(resp.data.yesterday);
        $('#overviewLastSevenDays').text(resp.data.last_seven_days);
        $('#overviewLastThirtyDays').text(resp.data.last_thirty_days);
        $('#overviewTotal').text(resp.data.total);
    }, 'json');
}

function loadInfo() {
    goGetChart.reset();
    $.get("/report/info", $('#reportFilterForm').serialize(), function(resp) {
        var lines = resp.data
        lineChartData = {
            labels: [],
        };
        dataset = [];
        for (var i = 0; i < lines.length; i++) {
            lineChartData.labels.push(lines[i].date.substring(0, 10));
            dataset.push(lines[i].count);
        }
        lineChartData.datasets = [
            {
                label: "go-get",
                // lineTension: 0,
                data: dataset,
                // fill: false,
                backgroundColor: 'rgba(0, 113, 206, 0.75)',
                borderColor: 'rgb(87, 146, 221)',
            }
        ]; 

        goGetChart.data = lineChartData
        goGetChart.update();
    }, 'json');
}