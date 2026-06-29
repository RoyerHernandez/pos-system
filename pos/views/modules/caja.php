<?php

$cashRegister = CashRegisterModel::mdlGetOpenCashRegister("caja", $_SESSION["id"]);

?>

<div class="content-wrapper">

  <section class="content-header">
    <h1>
      Caja
      <small>Apertura y cierre de caja</small>
    </h1>
    <ol class="breadcrumb">
      <li><a href="inicio"><i class="fa fa-dashboard"></i> Inicio</a></li>
      <li class="active">Caja</li>
    </ol>
  </section>

  <section class="content">

    <div class="row">

      <!-- Cash Register Status -->
      <div class="col-md-6 col-md-offset-3">

        <?php if($cashRegister): ?>

        <!-- CASH REGISTER OPEN -->
        <div class="box box-success">
          <div class="box-header with-border text-center">
            <h3 class="box-title"><i class="fa fa-unlock text-success"></i> Caja Abierta - #<?php echo $cashRegister['id']; ?></h3>
          </div>
          <div class="box-body">

            <table class="table">
              <tr>
                <td><strong>Fecha Apertura:</strong></td>
                <td class="text-right"><?php echo date("d/m/Y H:i", strtotime($cashRegister['fecha_apertura'])); ?></td>
              </tr>
              <tr>
                <td><strong>Monto Apertura:</strong></td>
                <td class="text-right">$<?php echo number_format($cashRegister['monto_apertura'], 2); ?></td>
              </tr>
              <tr class="info">
                <td><strong>Total Ventas:</strong></td>
                <td class="text-right"><strong>$<?php echo number_format($cashRegister['total_ventas'], 2); ?></strong></td>
              </tr>
              <tr>
                <td><i class="fa fa-money text-success"></i> Efectivo:</td>
                <td class="text-right">$<?php echo number_format($cashRegister['total_efectivo'], 2); ?></td>
              </tr>
              <tr>
                <td><i class="fa fa-credit-card text-info"></i> Tarjeta:</td>
                <td class="text-right">$<?php echo number_format($cashRegister['total_tarjeta'], 2); ?></td>
              </tr>
              <tr>
                <td><i class="fa fa-exchange text-warning"></i> Transferencia:</td>
                <td class="text-right">$<?php echo number_format($cashRegister['total_transferencia'], 2); ?></td>
              </tr>
              <tr class="success" style="font-size:18px;">
                <td><strong>Saldo Esperado en Caja:</strong></td>
                <td class="text-right"><strong>$<?php echo number_format($cashRegister['monto_apertura'] + $cashRegister['total_efectivo'], 2); ?></strong></td>
              </tr>
            </table>

            <hr>

            <form id="formCloseCashRegister">
              <div class="form-group">
                <label>Monto de Cierre (efectivo contado) ($)</label>
                <input type="number" class="form-control input-lg" id="closingAmount" name="closingAmount" step="0.01" min="0" placeholder="0.00" required style="text-align:center; font-size:24px;">
              </div>
              <button type="submit" class="btn btn-danger btn-lg btn-block">
                <i class="fa fa-lock"></i> Cerrar Caja
              </button>
            </form>

          </div>
        </div>

        <div class="text-center">
          <a href="crear-venta" class="btn btn-success">
            <i class="fa fa-cart-plus"></i> Ir al Punto de Venta
          </a>
        </div>

        <?php else: ?>

        <!-- CASH REGISTER CLOSED -->
        <div class="box box-warning">
          <div class="box-header with-border text-center">
            <h3 class="box-title"><i class="fa fa-lock"></i> Caja Cerrada</h3>
          </div>
          <div class="box-body text-center">
            <p>No hay caja abierta. Abra una caja para comenzar a vender.</p>
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

        <?php endif; ?>

      </div>

    </div>

  </section>

</div>
