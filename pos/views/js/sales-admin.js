/*=============================================
SALES ADMIN - Sales management module
DataTable, date filter, view details, cancel sales
=============================================*/

/*=============================================
DATATABLE INIT
=============================================*/

$(document).ready(function(){

	$(".salesTable").DataTable({
		"responsive": true,
		"order": [[0, "desc"]],
		"language": {
			"sProcessing": "Procesando...",
			"sLengthMenu": "Mostrar _MENU_ registros",
			"sZeroRecords": "No se encontraron ventas",
			"sEmptyTable": "No hay ventas registradas",
			"sInfo": "Mostrando _START_ a _END_ de _TOTAL_ ventas",
			"sInfoEmpty": "Mostrando 0 a 0 de 0 ventas",
			"sInfoFiltered": "(filtrado de _MAX_ ventas totales)",
			"sSearch": "Buscar:",
			"oPaginate": {
				"sFirst": "Primero",
				"sLast": "Ultimo",
				"sNext": "Siguiente",
				"sPrevious": "Anterior"
			}
		}
	});

});

/*=============================================
FILTER SALES BY DATE
=============================================*/

$(document).on("submit", "#formFilterSales", function(e){

	e.preventDefault();

	var startDate = $("#filterStartDate").val();
	var endDate = $("#filterEndDate").val();

	if(!startDate || !endDate){
		swal({ type: "warning", title: "Seleccione ambas fechas", showConfirmButton: true });
		return;
	}

	if(startDate > endDate){
		swal({ type: "warning", title: "La fecha inicial no puede ser mayor a la final", showConfirmButton: true });
		return;
	}

	$.ajax({
		url: "ajax/sales.ajax.php",
		method: "POST",
		data: { getSales: true, startDate: startDate, endDate: endDate },
		dataType: "json",
		success: function(response){
			rebuildSalesTable(response);
		}
	});

});

/*=============================================
SHOW ALL SALES
=============================================*/

$(document).on("click", "#btnShowAllSales", function(){

	$.ajax({
		url: "ajax/sales.ajax.php",
		method: "POST",
		data: { getSales: true },
		dataType: "json",
		success: function(response){
			rebuildSalesTable(response);
		}
	});

});

/*=============================================
REBUILD SALES TABLE
=============================================*/

function rebuildSalesTable(sales){

	// Destroy existing DataTable
	if($.fn.DataTable.isDataTable(".salesTable")){
		$(".salesTable").DataTable().destroy();
	}

	var html = "";

	if(sales && sales.length > 0){

		for(var i = 0; i < sales.length; i++){

			var s = sales[i];

			// Payment method icon/color
			var paymentIcon = "fa-money";
			var paymentColor = "success";
			if(s.metodo_pago == "Tarjeta"){ paymentIcon = "fa-credit-card"; paymentColor = "info"; }
			if(s.metodo_pago == "Transferencia"){ paymentIcon = "fa-exchange"; paymentColor = "warning"; }

			// Status badge
			var statusClass = "success";
			var statusLabel = "Completada";
			if(s.estado == "cancelada"){ statusClass = "danger"; statusLabel = "Cancelada"; }
			if(s.estado == "pendiente"){ statusClass = "warning"; statusLabel = "Pendiente"; }

			// Cancel button
			var cancelBtn = "";
			if(s.estado == "completada"){
				cancelBtn = '<button class="btn btn-danger btn-xs btnCancelSale" data-id="' + s.id + '" data-code="' + s.codigo_venta + '" title="Cancelar venta"><i class="fa fa-ban"></i></button>';
			}

			// Format date
			var date = new Date(s.fecha);
			var formattedDate = ("0" + date.getDate()).slice(-2) + "/" + ("0" + (date.getMonth()+1)).slice(-2) + "/" + date.getFullYear() + " " + ("0" + date.getHours()).slice(-2) + ":" + ("0" + date.getMinutes()).slice(-2);

			html += '<tr>' +
				'<td>' + (i + 1) + '</td>' +
				'<td><strong>' + s.codigo_venta + '</strong></td>' +
				'<td>' + formattedDate + '</td>' +
				'<td>' + (s.vendedor || '') + '</td>' +
				'<td>' + (s.cliente || '') + '</td>' +
				'<td><span class="text-' + paymentColor + '"><i class="fa ' + paymentIcon + '"></i> ' + s.metodo_pago + '</span></td>' +
				'<td><strong>$' + parseFloat(s.total).toFixed(2) + '</strong></td>' +
				'<td><span class="label label-' + statusClass + '">' + statusLabel + '</span></td>' +
				'<td><div class="btn-group">' +
					'<button class="btn btn-info btn-xs btnViewSale" data-id="' + s.id + '" data-code="' + s.codigo_venta + '" title="Ver detalle"><i class="fa fa-eye"></i></button>' +
					cancelBtn +
				'</div></td>' +
				'</tr>';
		}

	}

	$(".salesTable tbody").html(html);

	// Reinitialize DataTable
	$(".salesTable").DataTable({
		"responsive": true,
		"order": [[0, "desc"]],
		"language": {
			"sProcessing": "Procesando...",
			"sLengthMenu": "Mostrar _MENU_ registros",
			"sZeroRecords": "No se encontraron ventas",
			"sEmptyTable": "No hay ventas registradas",
			"sInfo": "Mostrando _START_ a _END_ de _TOTAL_ ventas",
			"sInfoEmpty": "Mostrando 0 a 0 de 0 ventas",
			"sInfoFiltered": "(filtrado de _MAX_ ventas totales)",
			"sSearch": "Buscar:",
			"oPaginate": {
				"sFirst": "Primero",
				"sLast": "Ultimo",
				"sNext": "Siguiente",
				"sPrevious": "Anterior"
			}
		}
	});

}

/*=============================================
VIEW SALE DETAIL
=============================================*/

$(document).on("click", ".btnViewSale", function(){

	var idSale = $(this).data("id");
	var saleCode = $(this).data("code");

	// Get sale info from the table row
	var row = $(this).closest("tr");
	var saleDate = row.find("td:eq(2)").text();
	var saleSeller = row.find("td:eq(3)").text();
	var saleClient = row.find("td:eq(4)").text();
	var salePayment = row.find("td:eq(5)").text().trim();
	var saleTotal = row.find("td:eq(6)").text();

	$("#modalSaleCode").text(saleCode);
	$("#modalSaleDate").text(saleDate);
	$("#modalSaleSeller").text(saleSeller);
	$("#modalSaleClient").text(saleClient);
	$("#modalSalePayment").text(salePayment);
	$("#modalSaleTotal").text(saleTotal);

	// Get sale details via AJAX
	$.ajax({
		url: "ajax/sales.ajax.php",
		method: "POST",
		data: { getSaleDetails: true, idSale: idSale },
		dataType: "json",
		success: function(response){

			var html = "";
			var subtotal = 0;
			var tax = 0;

			if(response && response.length > 0){

				for(var i = 0; i < response.length; i++){

					var d = response[i];
					subtotal += parseFloat(d.subtotal);

					html += '<tr>' +
						'<td>' + (i + 1) + '</td>' +
						'<td>' + d.codigo + '</td>' +
						'<td>' + d.producto + '</td>' +
						'<td>$' + parseFloat(d.precio_unitario).toFixed(2) + '</td>' +
						'<td>' + d.cantidad + '</td>' +
						'<td>$' + parseFloat(d.subtotal).toFixed(2) + '</td>' +
						'</tr>';
				}

			}else{
				html = '<tr><td colspan="6" class="text-center text-muted">Sin detalle disponible</td></tr>';
			}

			$("#modalSaleProducts").html(html);
			$("#modalSaleSubtotal").text("$" + subtotal.toFixed(2));
			$("#modalSaleTax").text("$" + tax.toFixed(2));

			$("#modalSaleDetail").modal("show");

		}
	});

});

/*=============================================
CANCEL SALE
=============================================*/

$(document).on("click", ".btnCancelSale", function(){

	var idSale = $(this).data("id");
	var saleCode = $(this).data("code");

	swal({
		title: "Cancelar venta " + saleCode + "?",
		text: "Esta accion restaurara el stock de los productos y no se puede deshacer",
		type: "warning",
		showCancelButton: true,
		confirmButtonColor: "#d33",
		cancelButtonColor: "#3085d6",
		confirmButtonText: "Si, cancelar venta",
		cancelButtonText: "No, mantener"
	}).then(function(result){

		if(result.value){

			$.ajax({
				url: "ajax/sales.ajax.php",
				method: "POST",
				data: { cancelSale: true, idSale: idSale },
				dataType: "json",
				success: function(response){

					if(response.status == "ok"){

						swal({
							type: "success",
							title: "Venta cancelada",
							text: "El stock ha sido restaurado",
							showConfirmButton: true
						}).then(function(){
							location.reload();
						});

					}else{

						swal({ type: "error", title: "Error al cancelar la venta", showConfirmButton: true });

					}

				}
			});

		}

	});

});
