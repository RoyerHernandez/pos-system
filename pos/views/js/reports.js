/*=============================================
REPORTS MODULE
Sales reports with charts and export
=============================================*/

var reportChartDaily = null;
var reportChartTopProducts = null;
var reportChartPayment = null;

/*=============================================
AUTO-LOAD ON PAGE READY
=============================================*/

$(document).ready(function(){

	var startDate = $("#reportStartDate").val();
	var endDate = $("#reportEndDate").val();

	if(startDate && endDate){
		loadAllReports(startDate, endDate);
	}

});

/*=============================================
GENERATE REPORT
=============================================*/

$(document).on("submit", "#formReportDates", function(e){

	e.preventDefault();

	var startDate = $("#reportStartDate").val();
	var endDate = $("#reportEndDate").val();

	if(!startDate || !endDate){
		swal({ type: "warning", title: "Seleccione ambas fechas", showConfirmButton: true });
		return;
	}

	if(startDate > endDate){
		swal({ type: "warning", title: "La fecha inicial no puede ser mayor a la final", showConfirmButton: true });
		return;
	}

	// Load active tab report
	var activeTab = $(".nav-tabs .active a").attr("href");

	loadAllReports(startDate, endDate);

});

// Also load on tab change
$(document).on("shown.bs.tab", 'a[data-toggle="tab"]', function(){

	var startDate = $("#reportStartDate").val();
	var endDate = $("#reportEndDate").val();

	if(startDate && endDate){
		var tab = $(this).attr("href");
		loadTabReport(tab, startDate, endDate);
	}

});

function loadAllReports(startDate, endDate){

	loadSalesByDate(startDate, endDate);
	loadTopProducts(startDate, endDate);
	loadByPaymentMethod(startDate, endDate);
	loadCashRegister(startDate, endDate);

}

function loadTabReport(tab, startDate, endDate){

	if(tab == "#tabSalesByDate") loadSalesByDate(startDate, endDate);
	if(tab == "#tabTopProducts") loadTopProducts(startDate, endDate);
	if(tab == "#tabByPaymentMethod") loadByPaymentMethod(startDate, endDate);
	if(tab == "#tabCashRegister") loadCashRegister(startDate, endDate);

}

/*=============================================
TAB 1: SALES BY DATE
=============================================*/

function loadSalesByDate(startDate, endDate){

	$.ajax({
		url: "ajax/reports.ajax.php",
		method: "POST",
		data: { reportSalesByDate: true, startDate: startDate, endDate: endDate },
		dataType: "json",
		success: function(response){

			// Summary cards
			var s = response.summary;
			var avg = s.completed_count > 0 ? (parseFloat(s.completed_amount) / parseInt(s.completed_count)) : 0;

			$("#summaryCompleted").text(s.completed_count);
			$("#summaryCompletedAmount").text("$" + parseFloat(s.completed_amount).toLocaleString(undefined, {minimumFractionDigits:2, maximumFractionDigits:2}));
			$("#summaryCancelled").text(s.cancelled_count);
			$("#summaryTotal").text(s.total_count);
			$("#summaryAvg").text("$" + avg.toFixed(2));
			$("#salesSummaryCards").show();

			// Daily chart
			if(response.daily && response.daily.length > 0){

				var labels = [];
				var data = [];

				for(var i = 0; i < response.daily.length; i++){
					var d = new Date(response.daily[i].sale_date);
					labels.push(("0" + d.getDate()).slice(-2) + "/" + ("0" + (d.getMonth()+1)).slice(-2));
					data.push(parseFloat(response.daily[i].total_amount));
				}

				if(reportChartDaily) reportChartDaily.destroy();

				reportChartDaily = new Chart(document.getElementById("chartDailySales").getContext("2d"), {
					type: "line",
					data: {
						labels: labels,
						datasets: [{
							label: "Ventas ($)",
							data: data,
							backgroundColor: "rgba(60, 141, 188, 0.2)",
							borderColor: "rgba(60, 141, 188, 1)",
							borderWidth: 2,
							pointRadius: 4,
							fill: true
						}]
					},
					options: {
						responsive: true,
						maintainAspectRatio: false,
						scales: {
							yAxes: [{ ticks: { beginAtZero: true, callback: function(v){ return "$" + v.toLocaleString(); } } }]
						}
					}
				});

				$("#dailyChartContainer").show();

			}else{
				$("#dailyChartContainer").hide();
			}

			// Sales table
			var sales = response.sales;
			var html = "";

			if(sales && sales.length > 0){

				html += '<table class="table table-bordered table-striped reportTable"><thead><tr>';
				html += '<th>#</th><th>Codigo</th><th>Fecha</th><th>Vendedor</th><th>Cliente</th><th>Metodo Pago</th><th>Total</th><th>Estado</th>';
				html += '</tr></thead><tbody>';

				for(var i = 0; i < sales.length; i++){

					var sale = sales[i];
					var date = new Date(sale.fecha);
					var formattedDate = ("0" + date.getDate()).slice(-2) + "/" + ("0" + (date.getMonth()+1)).slice(-2) + "/" + date.getFullYear() + " " + ("0" + date.getHours()).slice(-2) + ":" + ("0" + date.getMinutes()).slice(-2);

					var statusClass = "success";
					var statusLabel = "Completada";
					if(sale.estado == "cancelada"){ statusClass = "danger"; statusLabel = "Cancelada"; }

					html += '<tr>';
					html += '<td>' + (i+1) + '</td>';
					html += '<td><strong>' + sale.codigo_venta + '</strong></td>';
					html += '<td>' + formattedDate + '</td>';
					html += '<td>' + (sale.vendedor || '') + '</td>';
					html += '<td>' + (sale.cliente || '') + '</td>';
					html += '<td>' + sale.metodo_pago + '</td>';
					html += '<td><strong>$' + parseFloat(sale.total).toFixed(2) + '</strong></td>';
					html += '<td><span class="label label-' + statusClass + '">' + statusLabel + '</span></td>';
					html += '</tr>';

				}

				html += '</tbody></table>';

			}else{
				html = '<div class="text-center text-muted" style="padding:30px;"><p>No hay ventas en este periodo</p></div>';
			}

			// Destroy existing DataTable if any
			if($.fn.DataTable.isDataTable(".reportTable")){
				$(".reportTable").DataTable().destroy();
			}

			$("#salesByDateTable").html(html);

			if(sales && sales.length > 0){
				$(".reportTable").DataTable({
					"responsive": true,
					"order": [[0, "asc"]],
					"dom": '<"row"<"col-sm-6"l><"col-sm-6"f>>t<"row"<"col-sm-5"i><"col-sm-7"p>>B',
					"buttons": ["print", "csv"],
					"language": {
						"sZeroRecords": "No se encontraron ventas",
						"sEmptyTable": "Sin datos",
						"sInfo": "Mostrando _START_ a _END_ de _TOTAL_",
						"sSearch": "Buscar:",
						"oPaginate": { "sNext": "Sig.", "sPrevious": "Ant." }
					}
				});
			}

		}
	});

}

/*=============================================
TAB 2: TOP PRODUCTS
=============================================*/

function loadTopProducts(startDate, endDate){

	$.ajax({
		url: "ajax/reports.ajax.php",
		method: "POST",
		data: { reportTopProducts: true, startDate: startDate, endDate: endDate, limit: 10 },
		dataType: "json",
		success: function(response){

			if(response && response.length > 0){

				// Chart
				var labels = [];
				var dataQty = [];
				var dataRevenue = [];
				var colors = ["#00a65a","#00c0ef","#f39c12","#605ca8","#dd4b39","#3c8dbc","#d2d6de","#001f3f","#39cccc","#ff851b"];

				for(var i = 0; i < response.length; i++){
					labels.push(response[i].descripcion);
					dataQty.push(parseInt(response[i].total_qty));
					dataRevenue.push(parseFloat(response[i].total_revenue));
				}

				if(reportChartTopProducts) reportChartTopProducts.destroy();

				reportChartTopProducts = new Chart(document.getElementById("chartTopProductsReport").getContext("2d"), {
					type: "horizontalBar",
					data: {
						labels: labels,
						datasets: [{
							label: "Unidades vendidas",
							data: dataQty,
							backgroundColor: colors.slice(0, response.length),
							borderWidth: 1
						}]
					},
					options: {
						responsive: true,
						maintainAspectRatio: false,
						scales: { xAxes: [{ ticks: { beginAtZero: true, stepSize: 1 } }] },
						legend: { display: false }
					}
				});

				$("#topProductsChart").show();

				// Table
				var html = '<table class="table table-bordered table-striped"><thead><tr>';
				html += '<th>#</th><th>Codigo</th><th>Producto</th><th>Unidades</th><th>Precio Prom.</th><th>Ingresos</th>';
				html += '</tr></thead><tbody>';

				for(var i = 0; i < response.length; i++){
					html += '<tr>';
					html += '<td>' + (i+1) + '</td>';
					html += '<td>' + response[i].codigo + '</td>';
					html += '<td>' + response[i].descripcion + '</td>';
					html += '<td><strong>' + response[i].total_qty + '</strong></td>';
					html += '<td>$' + parseFloat(response[i].avg_price).toFixed(2) + '</td>';
					html += '<td><strong>$' + parseFloat(response[i].total_revenue).toFixed(2) + '</strong></td>';
					html += '</tr>';
				}

				html += '</tbody></table>';
				$("#topProductsTable").html(html);

			}else{

				$("#topProductsChart").hide();
				$("#topProductsTable").html('<div class="text-center text-muted" style="padding:30px;"><p>No hay datos en este periodo</p></div>');

			}

		}
	});

}

/*=============================================
TAB 3: BY PAYMENT METHOD
=============================================*/

function loadByPaymentMethod(startDate, endDate){

	$.ajax({
		url: "ajax/reports.ajax.php",
		method: "POST",
		data: { reportByPaymentMethod: true, startDate: startDate, endDate: endDate },
		dataType: "json",
		success: function(response){

			if(response && response.length > 0){

				var labels = [];
				var data = [];
				var colorMap = { "Efectivo": "#00a65a", "Tarjeta": "#00c0ef", "Transferencia": "#f39c12", "Mixto": "#605ca8" };
				var colors = [];
				var grandTotal = 0;

				for(var i = 0; i < response.length; i++){
					labels.push(response[i].metodo_pago);
					data.push(parseFloat(response[i].total_amount));
					colors.push(colorMap[response[i].metodo_pago] || "#999");
					grandTotal += parseFloat(response[i].total_amount);
				}

				// Chart
				if(reportChartPayment) reportChartPayment.destroy();

				reportChartPayment = new Chart(document.getElementById("chartPaymentReport").getContext("2d"), {
					type: "doughnut",
					data: {
						labels: labels,
						datasets: [{
							data: data,
							backgroundColor: colors,
							borderWidth: 2
						}]
					},
					options: {
						responsive: true,
						maintainAspectRatio: false,
						tooltips: {
							callbacks: {
								label: function(t, d){
									var label = d.labels[t.index];
									var value = d.datasets[0].data[t.index];
									var pct = ((value / grandTotal) * 100).toFixed(1);
									return label + ": $" + value.toLocaleString(undefined, {minimumFractionDigits:2}) + " (" + pct + "%)";
								}
							}
						}
					}
				});

				$("#paymentChartContainer").show();

				// Table
				var html = '<table class="table table-bordered"><thead><tr>';
				html += '<th>Metodo de Pago</th><th>Transacciones</th><th>Monto Total</th><th>Porcentaje</th>';
				html += '</tr></thead><tbody>';

				var iconMap = { "Efectivo": "fa-money text-success", "Tarjeta": "fa-credit-card text-info", "Transferencia": "fa-exchange text-warning", "Mixto": "fa-random text-purple" };

				for(var i = 0; i < response.length; i++){
					var pct = grandTotal > 0 ? ((parseFloat(response[i].total_amount) / grandTotal) * 100).toFixed(1) : 0;
					var icon = iconMap[response[i].metodo_pago] || "fa-question";

					html += '<tr>';
					html += '<td><i class="fa ' + icon + '"></i> ' + response[i].metodo_pago + '</td>';
					html += '<td class="text-center">' + response[i].total_count + '</td>';
					html += '<td class="text-right"><strong>$' + parseFloat(response[i].total_amount).toFixed(2) + '</strong></td>';
					html += '<td class="text-center">' + pct + '%</td>';
					html += '</tr>';
				}

				html += '<tr class="success"><td><strong>TOTAL</strong></td><td></td>';
				html += '<td class="text-right"><strong>$' + grandTotal.toFixed(2) + '</strong></td><td class="text-center">100%</td></tr>';
				html += '</tbody></table>';

				$("#paymentMethodTable").html(html);

			}else{

				$("#paymentChartContainer").hide();
				$("#paymentMethodTable").html('<div class="text-center text-muted" style="padding:30px;"><p>No hay datos en este periodo</p></div>');

			}

		}
	});

}

/*=============================================
TAB 4: CASH REGISTER
=============================================*/

function loadCashRegister(startDate, endDate){

	$.ajax({
		url: "ajax/reports.ajax.php",
		method: "POST",
		data: { reportCashRegister: true, startDate: startDate, endDate: endDate },
		dataType: "json",
		success: function(response){

			if(response && response.length > 0){

				var html = '<table class="table table-bordered table-striped"><thead><tr>';
				html += '<th>Caja #</th><th>Usuario</th><th>Apertura</th><th>Cierre</th><th>Monto Apertura</th><th>Monto Cierre</th><th>Total Ventas</th><th>Efectivo</th><th>Tarjeta</th><th>Transferencia</th><th>Estado</th>';
				html += '</tr></thead><tbody>';

				for(var i = 0; i < response.length; i++){

					var r = response[i];

					var openDate = new Date(r.fecha_apertura);
					var openFormatted = ("0" + openDate.getDate()).slice(-2) + "/" + ("0" + (openDate.getMonth()+1)).slice(-2) + "/" + openDate.getFullYear() + " " + ("0" + openDate.getHours()).slice(-2) + ":" + ("0" + openDate.getMinutes()).slice(-2);

					var closeFormatted = "-";
					if(r.fecha_cierre){
						var closeDate = new Date(r.fecha_cierre);
						closeFormatted = ("0" + closeDate.getDate()).slice(-2) + "/" + ("0" + (closeDate.getMonth()+1)).slice(-2) + "/" + closeDate.getFullYear() + " " + ("0" + closeDate.getHours()).slice(-2) + ":" + ("0" + closeDate.getMinutes()).slice(-2);
					}

					var statusClass = r.estado == "abierta" ? "success" : "default";

					html += '<tr>';
					html += '<td>' + r.id + '</td>';
					html += '<td>' + (r.usuario || '') + '</td>';
					html += '<td><small>' + openFormatted + '</small></td>';
					html += '<td><small>' + closeFormatted + '</small></td>';
					html += '<td class="text-right">$' + parseFloat(r.monto_apertura).toFixed(2) + '</td>';
					html += '<td class="text-right">' + (r.monto_cierre ? '$' + parseFloat(r.monto_cierre).toFixed(2) : '-') + '</td>';
					html += '<td class="text-right"><strong>$' + parseFloat(r.total_ventas).toFixed(2) + '</strong></td>';
					html += '<td class="text-right">$' + parseFloat(r.total_efectivo).toFixed(2) + '</td>';
					html += '<td class="text-right">$' + parseFloat(r.total_tarjeta).toFixed(2) + '</td>';
					html += '<td class="text-right">$' + parseFloat(r.total_transferencia).toFixed(2) + '</td>';
					html += '<td><span class="label label-' + statusClass + '">' + r.estado + '</span></td>';
					html += '</tr>';

				}

				html += '</tbody></table>';
				$("#cashRegisterTable").html(html);

			}else{

				$("#cashRegisterTable").html('<div class="text-center text-muted" style="padding:30px;"><p>No hay registros de caja en este periodo</p></div>');

			}

		}
	});

}
