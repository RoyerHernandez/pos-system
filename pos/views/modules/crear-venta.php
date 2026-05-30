<?php

// Check if cash register is open
$cashRegister = CashRegisterModel::mdlGetOpenCashRegister("caja", $_SESSION["id"]);

?>

<div class="content-wrapper">

  <section class="content-header">

    <h1>
      Punto de Venta
      <small>Crear nueva venta</small>
    </h1>

    <ol class="breadcrumb">
      <li><a href="inicio"><i class="fa fa-dashboard"></i> Inicio</a></li>
      <li class="active">Crear venta</li>
    </ol>

  </section>

  <section class="content">

    <?php if(!$cashRegister): ?>

    <!--=============================================
    CASH REGISTER CLOSED - OPEN MODAL
    =============================================-->

    <div class="row">
      <div class="col-md-6 col-md-offset-3">
        <div class="box box-warning">
          <div class="box-header with-border text-center">
            <h3 class="box-title"><i class="fa fa-lock"></i> Caja Cerrada</h3>
          </div>
          <div class="box-body text-center">
            <p>Debe abrir la caja antes de realizar ventas.</p>
            <form id="formOpenCashRegister">
              <div class="form-group">
                <label>Monto de Apertura ($)</label>
                <input type="number" class="form-control" id="openingAmount" name="openingAmount" step="0.01" min="0" value="0.00" required style="text-align:center; font-size:24px; max-width:300px; margin:0 auto;">
              </div>
              <button type="submit" class="btn btn-success btn-lg">
                <i class="fa fa-unlock"></i> Abrir Caja
              </button>
            </form>
          </div>
        </div>
      </div>
    </div>

    <?php else: ?>

    <!--=============================================
    POS INTERFACE
    =============================================-->

    <input type="hidden" id="idCashRegister" value="<?php echo $cashRegister['id']; ?>">

    <div class="row">

      <!--=============================================
      LEFT COLUMN: SEARCH + CART
      =============================================-->

      <div class="col-md-7">

        <!-- Product Search -->
        <div class="box box-primary">
          <div class="box-header with-border">
            <h3 class="box-title"><i class="fa fa-search"></i> Buscar Producto</h3>
            <span class="pull-right badge bg-green">F2</span>
          </div>
          <div class="box-body">
            <div class="input-group">
              <input type="text" class="form-control input-lg" id="searchProduct" placeholder="Buscar por nombre, codigo o escanear codigo de barras..." autocomplete="off">
              <span class="input-group-btn">
                <button class="btn btn-primary btn-lg" type="button" id="btnSearchProduct">
                  <i class="fa fa-search"></i>
                </button>
              </span>
            </div>
            <!-- Autocomplete Results -->
            <div id="searchResults" class="list-group" style="position:absolute; z-index:1000; width:calc(100% - 30px); max-height:300px; overflow-y:auto; display:none;">
            </div>
          </div>
        </div>

        <!-- Cart Table -->
        <div class="box box-success">
          <div class="box-header with-border">
            <h3 class="box-title"><i class="fa fa-shopping-cart"></i> Carrito de Venta</h3>
            <span class="pull-right">
              <button class="btn btn-danger btn-xs" id="btnClearCart" title="Vaciar carrito">
                <i class="fa fa-trash"></i> Vaciar
              </button>
            </span>
          </div>
          <div class="box-body table-responsive no-padding">
            <table class="table table-hover table-striped" id="cartTable">
              <thead>
                <tr>
                  <th style="width:5%">#</th>
                  <th style="width:10%">Codigo</th>
                  <th style="width:30%">Producto</th>
                  <th style="width:15%">Precio</th>
                  <th style="width:15%">Cantidad</th>
                  <th style="width:15%">Subtotal</th>
                  <th style="width:10%">Acciones</th>
                </tr>
              </thead>
              <tbody id="cartBody">
                <tr id="emptyCartRow">
                  <td colspan="7" class="text-center text-muted" style="padding:30px;">
                    <i class="fa fa-shopping-cart fa-3x"></i>
                    <p style="margin-top:10px;">El carrito esta vacio. Busque un producto para comenzar.</p>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

      </div>

      <!--=============================================
      RIGHT COLUMN: CLIENT + TOTALS + PAYMENT
      =============================================-->

      <div class="col-md-5">

        <!-- Client Selection -->
        <div class="box box-info">
          <div class="box-header with-border">
            <h3 class="box-title"><i class="fa fa-user"></i> Cliente</h3>
          </div>
          <div class="box-body">
            <div class="input-group">
              <input type="text" class="form-control" id="searchClient" placeholder="Buscar cliente..." autocomplete="off" value="Publico General">
              <span class="input-group-btn">
                <button class="btn btn-default" type="button" id="btnResetClient" title="Publico General">
                  <i class="fa fa-refresh"></i>
                </button>
              </span>
            </div>
            <input type="hidden" id="selectedClientId" value="1">
            <div id="clientResults" class="list-group" style="position:absolute; z-index:1000; width:calc(100% - 30px); max-height:200px; overflow-y:auto; display:none;">
            </div>
            <div id="clientInfo" class="text-muted small" style="margin-top:5px;">
              <i class="fa fa-check-circle text-success"></i> Publico General
            </div>
          </div>
        </div>

        <!-- Totals -->
        <div class="box box-default">
          <div class="box-header with-border">
            <h3 class="box-title"><i class="fa fa-calculator"></i> Resumen</h3>
          </div>
          <div class="box-body">
            <table class="table" style="margin-bottom:0;">
              <tr>
                <td><strong>Subtotal:</strong></td>
                <td class="text-right" id="displaySubtotal">$0.00</td>
              </tr>
              <tr>
                <td><strong>Impuesto (0%):</strong></td>
                <td class="text-right" id="displayTax">$0.00</td>
              </tr>
              <tr>
                <td><strong>Descuento:</strong></td>
                <td class="text-right" id="displayDiscount">$0.00</td>
              </tr>
              <tr class="success" style="font-size:24px; font-weight:bold;">
                <td><strong>TOTAL:</strong></td>
                <td class="text-right" id="displayTotal">$0.00</td>
              </tr>
            </table>
          </div>
        </div>

        <!-- Payment Method -->
        <div class="box box-primary">
          <div class="box-header with-border">
            <h3 class="box-title"><i class="fa fa-credit-card"></i> Metodo de Pago</h3>
            <span class="pull-right badge bg-blue">F4</span>
          </div>
          <div class="box-body">
            <div class="btn-group btn-group-justified" id="paymentMethods" role="group">
              <div class="btn-group" role="group">
                <button type="button" class="btn btn-success active" data-method="Efectivo">
                  <i class="fa fa-money"></i><br>Efectivo
                </button>
              </div>
              <div class="btn-group" role="group">
                <button type="button" class="btn btn-default" data-method="Tarjeta">
                  <i class="fa fa-credit-card"></i><br>Tarjeta
                </button>
              </div>
              <div class="btn-group" role="group">
                <button type="button" class="btn btn-default" data-method="Transferencia">
                  <i class="fa fa-exchange"></i><br>Transferencia
                </button>
              </div>
            </div>

            <!-- Cash Payment Section -->
            <div id="cashPaymentSection" style="margin-top:15px;">
              <div class="form-group">
                <label>Monto Recibido ($)</label>
                <input type="number" class="form-control input-lg" id="amountReceived" step="0.01" min="0" placeholder="0.00" style="text-align:center; font-size:20px;">
              </div>
              <div class="alert alert-info text-center" id="changeDisplay" style="font-size:20px; font-weight:bold; display:none;">
                Cambio: $<span id="changeAmount">0.00</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Complete Sale Button -->
        <button class="btn btn-lg btn-block btn-success" id="btnCompleteSale" disabled style="font-size:22px; padding:15px;">
          <i class="fa fa-check-circle"></i> COMPLETAR VENTA <small>(F12)</small>
        </button>

        <!-- Cash Register Info -->
        <div class="box box-default" style="margin-top:15px;">
          <div class="box-body text-center">
            <small class="text-muted">
              <i class="fa fa-desktop"></i> Caja #<?php echo $cashRegister['id']; ?> |
              Apertura: $<?php echo number_format($cashRegister['monto_apertura'], 2); ?> |
              Ventas: $<span id="cashRegisterSales"><?php echo number_format($cashRegister['total_ventas'], 2); ?></span>
            </small>
            <br>
            <a href="caja" class="btn btn-xs btn-default" style="margin-top:5px;">
              <i class="fa fa-cog"></i> Gestionar Caja
            </a>
          </div>
        </div>

      </div>

    </div>

    <?php endif; ?>

  </section>

</div>
