/*=============================================
POS - SALES MODULE
Cart management, product search, payment processing
=============================================*/

// Cart array
var cart = [];
var selectedPaymentMethod = "Efectivo";
var taxRate = 0; // 0% tax by default

/*=============================================
OPEN CASH REGISTER
=============================================*/

$(document).on("submit", "#formOpenCashRegister", function(e){

	e.preventDefault();

	var openingAmount = $("#openingAmount").val();

	$.ajax({
		url: "ajax/sales.ajax.php",
		method: "POST",
		data: { openCashRegister: true, openingAmount: openingAmount },
		dataType: "json",
		success: function(response){
			if(response.status == "ok"){
				swal({
					type: "success",
					title: "Caja abierta exitosamente",
					showConfirmButton: true,
					confirmButtonText: "Cerrar"
				}).then(function(result){
					window.location.reload();
				});
			}else{
				swal({ type: "error", title: "Error al abrir la caja", showConfirmButton: true });
			}
		}
	});

});

/*=============================================
PRODUCT SEARCH - AUTOCOMPLETE
=============================================*/

var searchTimeout;

$(document).on("keyup", "#searchProduct", function(){

	var search = $(this).val();

	clearTimeout(searchTimeout);

	if(search.length < 2){
		$("#searchResults").hide().empty();
		return;
	}

	searchTimeout = setTimeout(function(){

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
						var disabled = response[i].stock <= 0 ? "disabled" : "";

						html += '<a href="#" class="list-group-item searchResultItem ' + disabled + '" ' +
							'data-id="' + response[i].id + '" ' +
							'data-code="' + response[i].codigo + '" ' +
							'data-name="' + response[i].descripcion + '" ' +
							'data-price="' + response[i].precio_venta + '" ' +
							'data-stock="' + response[i].stock + '">' +
							'<strong>' + response[i].codigo + '</strong> - ' + response[i].descripcion +
							' <span class="pull-right ' + stockClass + '">$' + parseFloat(response[i].precio_venta).toFixed(2) +
							' | Stock: ' + response[i].stock + '</span>' +
							'</a>';
					}

				}else{
					html = '<a href="#" class="list-group-item disabled">No se encontraron productos</a>';
				}

				$("#searchResults").html(html).show();

			}
		});

	}, 300);

});

// Handle barcode scanner (Enter key in search)
$(document).on("keydown", "#searchProduct", function(e){

	if(e.keyCode == 13){ // Enter

		e.preventDefault();

		var code = $(this).val().trim();

		if(code.length > 0){

			$.ajax({
				url: "ajax/productos.ajax.php",
				method: "POST",
				data: { productCode: code },
				dataType: "json",
				success: function(response){

					if(response && response.id){

						addToCart(response.id, response.codigo, response.descripcion, parseFloat(response.precio_venta), parseInt(response.stock));
						$("#searchProduct").val("").focus();
						$("#searchResults").hide().empty();

					}else{
						swal({ type: "warning", title: "Producto no encontrado", timer: 1500, showConfirmButton: false });
					}
				}
			});
		}
	}

});

// Click on search result
$(document).on("click", ".searchResultItem", function(e){

	e.preventDefault();

	if($(this).hasClass("disabled")) return;

	var id = $(this).data("id");
	var code = $(this).data("code");
	var name = $(this).data("name");
	var price = parseFloat($(this).data("price"));
	var stock = parseInt($(this).data("stock"));

	addToCart(id, code, name, price, stock);

	$("#searchProduct").val("").focus();
	$("#searchResults").hide().empty();

});

// Hide search results on click outside
$(document).on("click", function(e){
	if(!$(e.target).closest("#searchProduct, #searchResults").length){
		$("#searchResults").hide();
	}
});

/*=============================================
CART MANAGEMENT
=============================================*/

function addToCart(id, code, name, price, stock){

	// Check if product already in cart
	var existing = cart.findIndex(function(item){ return item.id == id; });

	if(existing >= 0){

		// Check stock
		if(cart[existing].quantity + 1 > stock){
			swal({ type: "warning", title: "Stock insuficiente", text: "Solo hay " + stock + " unidades disponibles", timer: 2000, showConfirmButton: false });
			return;
		}

		cart[existing].quantity++;
		cart[existing].subtotal = cart[existing].quantity * cart[existing].price;

	}else{

		if(stock <= 0){
			swal({ type: "warning", title: "Producto sin stock", timer: 1500, showConfirmButton: false });
			return;
		}

		cart.push({
			id: id,
			code: code,
			name: name,
			price: price,
			quantity: 1,
			discount: 0,
			subtotal: price,
			stock: stock
		});

	}

	renderCart();

}

function renderCart(){

	var html = "";

	if(cart.length == 0){

		html = '<tr id="emptyCartRow">' +
			'<td colspan="7" class="text-center text-muted" style="padding:30px;">' +
			'<i class="fa fa-shopping-cart fa-3x"></i>' +
			'<p style="margin-top:10px;">El carrito esta vacio. Busque un producto para comenzar.</p>' +
			'</td></tr>';

		$("#btnCompleteSale").prop("disabled", true);

	}else{

		for(var i = 0; i < cart.length; i++){

			html += '<tr>' +
				'<td>' + (i + 1) + '</td>' +
				'<td><small>' + cart[i].code + '</small></td>' +
				'<td>' + cart[i].name + '</td>' +
				'<td>$' + cart[i].price.toFixed(2) + '</td>' +
				'<td>' +
					'<div class="input-group input-group-sm" style="width:120px;">' +
						'<span class="input-group-btn">' +
							'<button class="btn btn-danger btnDecreaseQty" data-index="' + i + '"><i class="fa fa-minus"></i></button>' +
						'</span>' +
						'<input type="number" class="form-control text-center inputQty" data-index="' + i + '" value="' + cart[i].quantity + '" min="1" max="' + cart[i].stock + '" style="width:50px;">' +
						'<span class="input-group-btn">' +
							'<button class="btn btn-success btnIncreaseQty" data-index="' + i + '"><i class="fa fa-plus"></i></button>' +
						'</span>' +
					'</div>' +
				'</td>' +
				'<td><strong>$' + cart[i].subtotal.toFixed(2) + '</strong></td>' +
				'<td>' +
					'<button class="btn btn-danger btn-xs btnRemoveItem" data-index="' + i + '" title="Eliminar">' +
						'<i class="fa fa-trash"></i>' +
					'</button>' +
				'</td>' +
				'</tr>';
		}

		$("#btnCompleteSale").prop("disabled", false);

	}

	$("#cartBody").html(html);
	calculateTotals();

}

function calculateTotals(){

	var subtotal = 0;

	for(var i = 0; i < cart.length; i++){
		subtotal += cart[i].subtotal;
	}

	var tax = subtotal * (taxRate / 100);
	var total = subtotal + tax;

	$("#displaySubtotal").text("$" + subtotal.toFixed(2));
	$("#displayTax").text("$" + tax.toFixed(2));
	$("#displayTotal").text("$" + total.toFixed(2));

	// Update change if cash payment
	calculateChange();

}

function calculateChange(){

	var total = parseFloat($("#displayTotal").text().replace("$", "")) || 0;
	var received = parseFloat($("#amountReceived").val()) || 0;

	if(received > 0 && received >= total){
		var change = received - total;
		$("#changeAmount").text(change.toFixed(2));
		$("#changeDisplay").show();
	}else{
		$("#changeDisplay").hide();
	}

}

// Increase quantity
$(document).on("click", ".btnIncreaseQty", function(){
	var index = $(this).data("index");
	if(cart[index].quantity + 1 <= cart[index].stock){
		cart[index].quantity++;
		cart[index].subtotal = cart[index].quantity * cart[index].price;
		renderCart();
	}else{
		swal({ type: "warning", title: "Stock insuficiente", timer: 1500, showConfirmButton: false });
	}
});

// Decrease quantity
$(document).on("click", ".btnDecreaseQty", function(){
	var index = $(this).data("index");
	if(cart[index].quantity > 1){
		cart[index].quantity--;
		cart[index].subtotal = cart[index].quantity * cart[index].price;
		renderCart();
	}
});

// Manual quantity input
$(document).on("change", ".inputQty", function(){
	var index = $(this).data("index");
	var qty = parseInt($(this).val()) || 1;

	if(qty < 1) qty = 1;
	if(qty > cart[index].stock){
		qty = cart[index].stock;
		swal({ type: "warning", title: "Stock insuficiente", text: "Maximo: " + cart[index].stock, timer: 1500, showConfirmButton: false });
	}

	cart[index].quantity = qty;
	cart[index].subtotal = qty * cart[index].price;
	renderCart();
});

// Remove item
$(document).on("click", ".btnRemoveItem", function(){
	var index = $(this).data("index");
	cart.splice(index, 1);
	renderCart();
});

// Clear cart
$(document).on("click", "#btnClearCart", function(){

	if(cart.length == 0) return;

	swal({
		title: "Vaciar carrito?",
		text: "Se eliminaran todos los productos del carrito",
		type: "warning",
		showCancelButton: true,
		confirmButtonColor: "#d33",
		cancelButtonColor: "#3085d6",
		confirmButtonText: "Si, vaciar",
		cancelButtonText: "Cancelar"
	}).then(function(result){
		if(result.value){
			cart = [];
			renderCart();
		}
	});

});

/*=============================================
CLIENT SEARCH
=============================================*/

var clientTimeout;

$(document).on("keyup", "#searchClient", function(){

	var search = $(this).val();

	clearTimeout(clientTimeout);

	if(search.length < 2){
		$("#clientResults").hide().empty();
		return;
	}

	clientTimeout = setTimeout(function(){

		$.ajax({
			url: "ajax/clientes.ajax.php",
			method: "POST",
			data: { searchClient: search },
			dataType: "json",
			success: function(response){

				var html = "";

				if(response.length > 0){

					for(var i = 0; i < response.length; i++){

						html += '<a href="#" class="list-group-item clientResultItem" ' +
							'data-id="' + response[i].id + '" ' +
							'data-name="' + response[i].nombre + '">' +
							'<strong>' + response[i].nombre + '</strong>' +
							(response[i].documento ? ' <small class="text-muted">(' + response[i].documento + ')</small>' : '') +
							'</a>';

					}

				}else{
					html = '<a href="#" class="list-group-item disabled">No se encontraron clientes</a>';
				}

				$("#clientResults").html(html).show();

			}
		});

	}, 300);

});

// Click on client result
$(document).on("click", ".clientResultItem", function(e){

	e.preventDefault();

	var id = $(this).data("id");
	var name = $(this).data("name");

	$("#selectedClientId").val(id);
	$("#searchClient").val(name);
	$("#clientInfo").html('<i class="fa fa-check-circle text-success"></i> ' + name);
	$("#clientResults").hide().empty();

});

// Reset to Publico General
$(document).on("click", "#btnResetClient", function(){
	$("#selectedClientId").val(1);
	$("#searchClient").val("Publico General");
	$("#clientInfo").html('<i class="fa fa-check-circle text-success"></i> Publico General');
});

// Hide client results on click outside
$(document).on("click", function(e){
	if(!$(e.target).closest("#searchClient, #clientResults").length){
		$("#clientResults").hide();
	}
});

/*=============================================
PAYMENT METHOD SELECTION
=============================================*/

$(document).on("click", "#paymentMethods button", function(){

	$("#paymentMethods button").removeClass("active btn-success btn-info btn-warning").addClass("btn-default");
	$(this).removeClass("btn-default").addClass("active");

	selectedPaymentMethod = $(this).data("method");

	if(selectedPaymentMethod == "Efectivo"){
		$(this).addClass("btn-success");
		$("#cashPaymentSection").show();
	}else if(selectedPaymentMethod == "Tarjeta"){
		$(this).addClass("btn-info");
		$("#cashPaymentSection").hide();
	}else{
		$(this).addClass("btn-warning");
		$("#cashPaymentSection").hide();
	}

});

// Calculate change on input
$(document).on("keyup", "#amountReceived", function(){
	calculateChange();
});

/*=============================================
COMPLETE SALE
=============================================*/

$(document).on("click", "#btnCompleteSale", function(){

	if(cart.length == 0){
		swal({ type: "warning", title: "El carrito esta vacio", showConfirmButton: true });
		return;
	}

	var subtotal = 0;
	for(var i = 0; i < cart.length; i++){
		subtotal += cart[i].subtotal;
	}
	var tax = subtotal * (taxRate / 100);
	var total = subtotal + tax;

	// Validate cash payment
	if(selectedPaymentMethod == "Efectivo"){
		var received = parseFloat($("#amountReceived").val()) || 0;
		if(received < total){
			swal({ type: "warning", title: "Monto insuficiente", text: "El monto recibido debe ser mayor o igual al total", showConfirmButton: true });
			$("#amountReceived").focus();
			return;
		}
	}

	// Prepare products JSON
	var products = JSON.stringify(cart);

	swal({
		title: "Completar venta?",
		html: "<strong>Total: $" + total.toFixed(2) + "</strong><br>Metodo: " + selectedPaymentMethod,
		type: "question",
		showCancelButton: true,
		confirmButtonColor: "#28a745",
		cancelButtonColor: "#6c757d",
		confirmButtonText: "Si, completar venta",
		cancelButtonText: "Cancelar"
	}).then(function(result){

		if(result.value){

			$.ajax({
				url: "ajax/sales.ajax.php",
				method: "POST",
				data: {
					completeSale: true,
					idClient: $("#selectedClientId").val(),
					products: products,
					subtotal: subtotal.toFixed(2),
					tax: tax.toFixed(2),
					discount: "0.00",
					total: total.toFixed(2),
					paymentMethod: selectedPaymentMethod
				},
				dataType: "json",
				success: function(response){

					if(response.status == "ok"){

						var changeHtml = "";
						if(selectedPaymentMethod == "Efectivo"){
							var received = parseFloat($("#amountReceived").val()) || 0;
							var change = received - total;
							if(change > 0){
								changeHtml = "<br><br><strong style='font-size:24px; color:#28a745;'>Cambio: $" + change.toFixed(2) + "</strong>";
							}
						}

						swal({
							type: "success",
							title: "Venta completada!",
							html: "Codigo: <strong>" + response.saleCode + "</strong>" + changeHtml,
							showConfirmButton: true,
							confirmButtonText: "Nueva Venta"
						}).then(function(){
							// Reset POS
							cart = [];
							renderCart();
							$("#selectedClientId").val(1);
							$("#searchClient").val("Publico General");
							$("#clientInfo").html('<i class="fa fa-check-circle text-success"></i> Publico General');
							$("#amountReceived").val("");
							$("#changeDisplay").hide();
							$("#searchProduct").focus();

							// Update cash register sales display
							var currentSales = parseFloat($("#cashRegisterSales").text().replace(",", "")) || 0;
							$("#cashRegisterSales").text((currentSales + total).toFixed(2));
						});

					}else{

						swal({ type: "error", title: "Error", text: response.message || "Error al procesar la venta", showConfirmButton: true });

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
KEYBOARD SHORTCUTS
=============================================*/

$(document).on("keydown", function(e){

	// F2 - Focus search product
	if(e.keyCode == 113){
		e.preventDefault();
		$("#searchProduct").focus();
	}

	// F4 - Focus payment
	if(e.keyCode == 115){
		e.preventDefault();
		if(selectedPaymentMethod == "Efectivo"){
			$("#amountReceived").focus();
		}else{
			$("#paymentMethods button:first").focus();
		}
	}

	// F12 - Complete sale
	if(e.keyCode == 123){
		e.preventDefault();
		$("#btnCompleteSale").click();
	}

});

/*=============================================
CLOSE CASH REGISTER (from caja page)
=============================================*/

$(document).on("submit", "#formCloseCashRegister", function(e){

	e.preventDefault();

	var closingAmount = $("#closingAmount").val();

	swal({
		title: "Cerrar caja?",
		text: "Esta accion cerrara la caja actual",
		type: "warning",
		showCancelButton: true,
		confirmButtonColor: "#d33",
		confirmButtonText: "Si, cerrar caja",
		cancelButtonText: "Cancelar"
	}).then(function(result){

		if(result.value){

			$.ajax({
				url: "ajax/sales.ajax.php",
				method: "POST",
				data: { closeCashRegister: true, closingAmount: closingAmount },
				dataType: "json",
				success: function(response){
					if(response.status == "ok"){

						var cr = response.cashRegister;

						swal({
							type: "success",
							title: "Caja cerrada",
							html: "<table class='table table-condensed text-left'>" +
								"<tr><td>Total Ventas:</td><td class='text-right'>$" + parseFloat(cr.total_ventas).toFixed(2) + "</td></tr>" +
								"<tr><td>Efectivo:</td><td class='text-right'>$" + parseFloat(cr.total_efectivo).toFixed(2) + "</td></tr>" +
								"<tr><td>Tarjeta:</td><td class='text-right'>$" + parseFloat(cr.total_tarjeta).toFixed(2) + "</td></tr>" +
								"<tr><td>Transferencia:</td><td class='text-right'>$" + parseFloat(cr.total_transferencia).toFixed(2) + "</td></tr>" +
								"</table>",
							showConfirmButton: true
						}).then(function(){
							window.location.reload();
						});

					}else{
						swal({ type: "error", title: "Error al cerrar la caja", showConfirmButton: true });
					}
				}
			});

		}

	});

});
