/*=============================================
INVENTORY MODULE
Movement history, manual entries/exits
=============================================*/

// Fixed catalog (must match InventoryController)
var entryMotivos = ["Compra", "Devolución de cliente", "Ajuste positivo"];
var exitMotivos = ["Merma", "Devolución a proveedor", "Ajuste negativo"];

/*=============================================
LOAD MOVEMENTS
=============================================*/

var movementsTable = null;

function loadMovements(startDate, endDate){

	var data = { getMovements: true };

	if(startDate && endDate){
		data.startDate = startDate;
		data.endDate = endDate;
	}

	$.ajax({
		url: "ajax/inventario.ajax.php",
		method: "POST",
		data: data,
		dataType: "json",
		success: function(response){

			var html = "";

			if(response && response.length > 0){

				for(var i = 0; i < response.length; i++){

					var mov = response[i];
					var tipoBadge = mov.tipo == "entrada"
						? '<span class="label label-success">Entrada</span>'
						: '<span class="label label-danger">Salida</span>';

					var fecha = mov.fecha ? mov.fecha.substring(0, 16).replace("T", " ") : "";

					html += '<tr>' +
						'<td>' + mov.id + '</td>' +
						'<td>' + fecha + '</td>' +
						'<td><strong>' + (mov.codigo_producto || '') + '</strong> - ' + (mov.producto || 'N/A') + '</td>' +
						'<td>' + tipoBadge + '</td>' +
						'<td>' + mov.motivo + '</td>' +
						'<td class="text-center"><strong>' + mov.cantidad + '</strong></td>' +
						'<td>' + (mov.usuario || 'N/A') + '</td>' +
						'<td>' + (mov.observaciones || '-') + '</td>' +
						'</tr>';
				}

			}

			$("#movementsBody").html(html);

			// Destroy existing DataTable if initialized
			if(movementsTable){
				movementsTable.destroy();
			}

			// Re-initialize DataTable
			movementsTable = $("#tableMovements").DataTable({
				"order": [[0, "desc"]],
				"language": {
					"sProcessing": "Procesando...",
					"sLengthMenu": "Mostrar _MENU_ registros",
					"sZeroRecords": "No se encontraron resultados",
					"sEmptyTable": "Ningún dato disponible en esta tabla",
					"sInfo": "Mostrando registros del _START_ al _END_ de un total de _TOTAL_ registros",
					"sInfoEmpty": "Mostrando registros del 0 al 0 de un total de 0 registros",
					"sInfoFiltered": "(filtrado de un total de _MAX_ registros)",
					"sSearch": "Buscar:",
					"oPaginate": {
						"sFirst": "Primero",
						"sPrevious": "Anterior",
						"sNext": "Siguiente",
						"sLast": "Ultimo"
					}
				}
			});

		}
	});

}

// Load all movements on page load (only if on inventory page)
$(document).ready(function(){
	if($("#tableMovements").length > 0){
		loadMovements(null, null);
	}
});

/*=============================================
DATE RANGE FILTER
=============================================*/

$(document).on("submit", "#formFilterMovements", function(e){

	e.preventDefault();

	var startDate = $("#movStartDate").val();
	var endDate = $("#movEndDate").val();

	if(startDate && endDate){
		loadMovements(startDate, endDate);
	}

});

$(document).on("click", "#btnClearFilter", function(){

	var now = new Date();
	var year = now.getFullYear();
	var month = ("0" + (now.getMonth() + 1)).slice(-2);
	var day = ("0" + now.getDate()).slice(-2);

	$("#movStartDate").val(year + "-" + month + "-01");
	$("#movEndDate").val(year + "-" + month + "-" + day);
	loadMovements(null, null);

});

/*=============================================
TYPE CHANGE - UPDATE MOTIVO DROPDOWN
=============================================*/

$(document).on("change", "#movType", function(){

	var tipo = $(this).val();
	var motivos = [];

	if(tipo == "entrada"){
		motivos = entryMotivos;
	}else if(tipo == "salida"){
		motivos = exitMotivos;
	}

	var html = '<option value="">-- Seleccionar motivo --</option>';

	for(var i = 0; i < motivos.length; i++){
		html += '<option value="' + motivos[i] + '">' + motivos[i] + '</option>';
	}

	$("#movMotivo").html(html).prop("disabled", motivos.length == 0);

});

/*=============================================
PRODUCT SEARCH IN MODAL
=============================================*/

var movProductTimeout;

$(document).on("keyup", "#searchMovProduct", function(){

	var search = $(this).val();

	clearTimeout(movProductTimeout);

	if(search.length < 2){
		$("#movProductResults").hide().empty();
		return;
	}

	movProductTimeout = setTimeout(function(){

		$.ajax({
			url: "ajax/productos.ajax.php",
			method: "POST",
			data: { searchProduct: search },
			dataType: "json",
			success: function(response){

				var html = "";

				if(response.length > 0){

					for(var i = 0; i < response.length; i++){

						var stockClass = response[i].stock <= 0 ? "text-danger" : "text-success";

						html += '<a href="#" class="list-group-item movProductItem" ' +
							'data-id="' + response[i].id + '" ' +
							'data-code="' + response[i].codigo + '" ' +
							'data-name="' + response[i].descripcion + '" ' +
							'data-stock="' + response[i].stock + '">' +
							'<strong>' + response[i].codigo + '</strong> - ' + response[i].descripcion +
							' <span class="pull-right ' + stockClass + '">Stock: ' + response[i].stock + '</span>' +
							'</a>';
					}

				}else{
					html = '<a href="#" class="list-group-item disabled">No se encontraron productos</a>';
				}

				$("#movProductResults").html(html).show();

			}
		});

	}, 300);

});

// Select product from results
$(document).on("click", ".movProductItem", function(e){

	e.preventDefault();

	var id = $(this).data("id");
	var code = $(this).data("code");
	var name = $(this).data("name");
	var stock = $(this).data("stock");

	$("#movProductId").val(id);
	$("#searchMovProduct").val(code + " - " + name);
	$("#movProductInfo").html('<i class="fa fa-check-circle text-success"></i> ' + code + ' - ' + name + ' | Stock: ' + stock);
	$("#movProductResults").hide().empty();

});

// Hide product results on click outside
$(document).on("click", function(e){
	if(!$(e.target).closest("#searchMovProduct, #movProductResults").length){
		$("#movProductResults").hide();
	}
});

/*=============================================
SUBMIT NEW MOVEMENT
=============================================*/

$(document).on("submit", "#formNewMovement", function(e){

	e.preventDefault();

	var idProducto = $("#movProductId").val();
	var tipo = $("#movType").val();
	var motivo = $("#movMotivo").val();
	var cantidad = $("#movCantidad").val();
	var observaciones = $("#movObservaciones").val();

	// Validate
	if(!idProducto || idProducto == ""){
		swal({ type: "warning", title: "Seleccione un producto", showConfirmButton: true });
		return;
	}

	if(!tipo || tipo == ""){
		swal({ type: "warning", title: "Seleccione un tipo de movimiento", showConfirmButton: true });
		return;
	}

	if(!motivo || motivo == ""){
		swal({ type: "warning", title: "Seleccione un motivo", showConfirmButton: true });
		return;
	}

	if(!cantidad || parseInt(cantidad) <= 0){
		swal({ type: "warning", title: "La cantidad debe ser mayor a 0", showConfirmButton: true });
		return;
	}

	var tipoLabel = tipo == "entrada" ? "ENTRADA" : "SALIDA";

	swal({
		title: "Registrar movimiento?",
		html: "<strong>" + tipoLabel + "</strong>: " + motivo + "<br>Cantidad: <strong>" + cantidad + "</strong>",
		type: "question",
		showCancelButton: true,
		confirmButtonColor: "#3c8dbc",
		cancelButtonColor: "#6c757d",
		confirmButtonText: "Si, registrar",
		cancelButtonText: "Cancelar"
	}).then(function(result){

		if(result.value){

			$.ajax({
				url: "ajax/inventario.ajax.php",
				method: "POST",
				data: {
					createMovement: true,
					idProducto: idProducto,
					tipo: tipo,
					motivo: motivo,
					cantidad: cantidad,
					observaciones: observaciones
				},
				dataType: "json",
				success: function(response){

					if(response.status == "ok"){

						swal({
							type: "success",
							title: response.message,
							showConfirmButton: true,
							confirmButtonText: "Cerrar"
						}).then(function(){
							// Close modal and reset form
							$("#modalNewMovement").modal("hide");
							$("#formNewMovement")[0].reset();
							$("#movProductId").val("");
							$("#movProductInfo").html("");
							$("#movMotivo").html('<option value="">-- Seleccione tipo primero --</option>').prop("disabled", true);

							// Reload movements
							loadMovements(null, null);
						});

					}else{

						swal({ type: "error", title: "Error", text: response.message || "Error al registrar el movimiento", showConfirmButton: true });

					}

				},
				error: function(){
					swal({ type: "error", title: "Error de conexion", showConfirmButton: true });
				}
			});

		}

	});

});

/*=============================================
RESET MODAL ON CLOSE
=============================================*/

$(document).on("hidden.bs.modal", "#modalNewMovement", function(){
	$("#formNewMovement")[0].reset();
	$("#movProductId").val("");
	$("#movProductInfo").html("");
	$("#movMotivo").html('<option value="">-- Seleccione tipo primero --</option>').prop("disabled", true);
	$("#movProductResults").hide().empty();
});
