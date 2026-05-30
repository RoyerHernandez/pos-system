/*=============================================
DASHBOARD - Charts initialization
=============================================*/

$(document).ready(function(){

	if(typeof dashboardData === "undefined") return;

	/*=============================================
	SALES LAST 7 DAYS - BAR CHART
	=============================================*/

	var ctx7Days = document.getElementById("chartSales7Days");

	if(ctx7Days){

		new Chart(ctx7Days.getContext("2d"), {
			type: "bar",
			data: {
				labels: dashboardData.sales7Days.labels,
				datasets: [{
					label: "Ventas ($)",
					data: dashboardData.sales7Days.data,
					backgroundColor: "rgba(60, 141, 188, 0.7)",
					borderColor: "rgba(60, 141, 188, 1)",
					borderWidth: 1
				}]
			},
			options: {
				responsive: true,
				maintainAspectRatio: false,
				scales: {
					yAxes: [{
						ticks: {
							beginAtZero: true,
							callback: function(value){
								return "$" + value.toLocaleString();
							}
						}
					}]
				},
				tooltips: {
					callbacks: {
						label: function(tooltipItem){
							return "Ventas: $" + tooltipItem.yLabel.toLocaleString(undefined, {minimumFractionDigits: 2});
						}
					}
				}
			}
		});

	}

	/*=============================================
	TOP 5 PRODUCTS - HORIZONTAL BAR CHART
	=============================================*/

	var ctxTop = document.getElementById("chartTopProducts");

	if(ctxTop){

		new Chart(ctxTop.getContext("2d"), {
			type: "horizontalBar",
			data: {
				labels: dashboardData.topProducts.labels,
				datasets: [{
					label: "Unidades vendidas",
					data: dashboardData.topProducts.data,
					backgroundColor: [
						"rgba(0, 166, 90, 0.7)",
						"rgba(0, 192, 239, 0.7)",
						"rgba(243, 156, 18, 0.7)",
						"rgba(96, 92, 168, 0.7)",
						"rgba(221, 75, 57, 0.7)"
					],
					borderColor: [
						"rgba(0, 166, 90, 1)",
						"rgba(0, 192, 239, 1)",
						"rgba(243, 156, 18, 1)",
						"rgba(96, 92, 168, 1)",
						"rgba(221, 75, 57, 1)"
					],
					borderWidth: 1
				}]
			},
			options: {
				responsive: true,
				maintainAspectRatio: false,
				scales: {
					xAxes: [{
						ticks: {
							beginAtZero: true,
							stepSize: 1
						}
					}]
				},
				legend: {
					display: false
				}
			}
		});

	}

	/*=============================================
	PAYMENT METHOD - DOUGHNUT CHART
	=============================================*/

	var ctxPayment = document.getElementById("chartPaymentMethod");

	if(ctxPayment){

		new Chart(ctxPayment.getContext("2d"), {
			type: "doughnut",
			data: {
				labels: dashboardData.paymentMethod.labels,
				datasets: [{
					data: dashboardData.paymentMethod.data,
					backgroundColor: dashboardData.paymentMethod.colors,
					borderWidth: 2
				}]
			},
			options: {
				responsive: true,
				maintainAspectRatio: false,
				tooltips: {
					callbacks: {
						label: function(tooltipItem, data){
							var label = data.labels[tooltipItem.index];
							var value = data.datasets[0].data[tooltipItem.index];
							return label + ": $" + value.toLocaleString(undefined, {minimumFractionDigits: 2});
						}
					}
				}
			}
		});

	}

});
