<?php

$sales = SaleController::ctrShowSales(null, null);

?>

<div class="content-wrapper">

  <section class="content-header">
    <h1>
      Administrar Ventas
      <small>Historial y gestion de ventas</small>
    </h1>
    <ol class="breadcrumb">
      <li><a href="inicio"><i class="fa fa-dashboard"></i> Inicio</a></li>
      <li class="active">Ventas</li>
    </ol>
  </section>

  <section class="content">

    <!-- Date Filter -->
    <div class="box box-info">
      <div class="box-header with-border">
        <h3 class="box-title"><i class="fa fa-filter"></i> Filtrar por Fecha</h3>
      </div>
      <div class="box-body">
        <form id="formFilterSales" class="form-inline">
          <div class="form-group">
            <label>Desde:</label>
            <input type="date" class="form-control" id="filterStartDate" value="<?php echo date('Y-m-d'); ?>">
          </div>
          <div class="form-group" style="margin-left:15px;">
            <label>Hasta:</label>
            <input type="date" class="form-control" id="filterEndDate" value="<?php echo date('Y-m-d'); ?>">
          </div>
          <button type="submit" class="btn btn-primary" style="margin-left:15px;">
            <i class="fa fa-search"></i> Filtrar
          </button>
          <button type="button" class="btn btn-default" id="btnShowAllSales" style="margin-left:5px;">
            <i class="fa fa-list"></i> Mostrar Todas
          </button>
        </form>
      </div>
    </div>

    <!-- Sales Table -->
    <div class="box box-primary">
      <div class="box-header with-border">
        <h3 class="box-title"><i class="fa fa-shopping-cart"></i> Listado de Ventas</h3>
        <div class="box-tools pull-right">
          <a href="crear-venta" class="btn btn-success btn-sm">
            <i class="fa fa-plus"></i> Nueva Venta
          </a>
        </div>
      </div>
      <div class="box-body">
        <table class="table table-bordered table-striped dt-responsive salesTable" width="100%">
          <thead>
            <tr>
              <th style="width:5%">#</th>
              <th style="width:12%">Codigo</th>
              <th style="width:15%">Fecha</th>
              <th style="width:15%">Vendedor</th>
              <th style="width:13%">Cliente</th>
              <th style="width:10%">Metodo Pago</th>
              <th style="width:10%">Total</th>
              <th style="width:10%">Estado</th>
              <th style="width:10%">Acciones</th>
            </tr>
          </thead>
          <tbody>

            <?php if($sales): ?>

              <?php $counter = 1; ?>

              <?php foreach($sales as $key => $value): ?>

              <tr>
                <td><?php echo $counter; ?></td>
                <td><strong><?php echo $value["codigo_venta"]; ?></strong></td>
                <td><?php echo date("d/m/Y H:i", strtotime($value["fecha"])); ?></td>
                <td><?php echo $value["vendedor"]; ?></td>
                <td><?php echo $value["cliente"]; ?></td>
                <td>
                  <?php
                    $paymentIcon = "fa-money";
                    $paymentColor = "success";
                    if($value["metodo_pago"] == "Tarjeta"){ $paymentIcon = "fa-credit-card"; $paymentColor = "info"; }
                    if($value["metodo_pago"] == "Transferencia"){ $paymentIcon = "fa-exchange"; $paymentColor = "warning"; }
                  ?>
                  <span class="text-<?php echo $paymentColor; ?>">
                    <i class="fa <?php echo $paymentIcon; ?>"></i> <?php echo $value["metodo_pago"]; ?>
                  </span>
                </td>
                <td><strong>$<?php echo number_format($value["total"], 2); ?></strong></td>
                <td>
                  <?php
                    $statusClass = "success";
                    $statusLabel = "Completada";
                    if($value["estado"] == "cancelada"){ $statusClass = "danger"; $statusLabel = "Cancelada"; }
                    if($value["estado"] == "pendiente"){ $statusClass = "warning"; $statusLabel = "Pendiente"; }
                  ?>
                  <span class="label label-<?php echo $statusClass; ?>"><?php echo $statusLabel; ?></span>
                </td>
                <td>
                  <div class="btn-group">
                    <button class="btn btn-info btn-xs btnViewSale" data-id="<?php echo $value["id"]; ?>" data-code="<?php echo $value["codigo_venta"]; ?>" title="Ver detalle">
                      <i class="fa fa-eye"></i>
                    </button>

                    <?php if($value["estado"] == "completada"): ?>
                    <button class="btn btn-danger btn-xs btnCancelSale" data-id="<?php echo $value["id"]; ?>" data-code="<?php echo $value["codigo_venta"]; ?>" title="Cancelar venta">
                      <i class="fa fa-ban"></i>
                    </button>
                    <?php endif; ?>
                  </div>
                </td>
              </tr>

              <?php $counter++; ?>

              <?php endforeach; ?>

            <?php endif; ?>

          </tbody>
        </table>
      </div>
    </div>

  </section>

</div>

<!--=============================================
MODAL: SALE DETAIL
=============================================-->

<div class="modal fade" id="modalSaleDetail" tabindex="-1" role="dialog">
  <div class="modal-dialog modal-lg" role="document">
    <div class="modal-content">
      <div class="modal-header bg-info">
        <button type="button" class="close" data-dismiss="modal"><span>&times;</span></button>
        <h4 class="modal-title"><i class="fa fa-file-text-o"></i> Detalle de Venta - <span id="modalSaleCode"></span></h4>
      </div>
      <div class="modal-body">

        <!-- Sale Info -->
        <div class="row">
          <div class="col-md-6">
            <p><strong>Fecha:</strong> <span id="modalSaleDate"></span></p>
            <p><strong>Vendedor:</strong> <span id="modalSaleSeller"></span></p>
          </div>
          <div class="col-md-6">
            <p><strong>Cliente:</strong> <span id="modalSaleClient"></span></p>
            <p><strong>Metodo de Pago:</strong> <span id="modalSalePayment"></span></p>
          </div>
        </div>

        <hr>

        <!-- Sale Products -->
        <table class="table table-bordered table-striped">
          <thead>
            <tr>
              <th>#</th>
              <th>Codigo</th>
              <th>Producto</th>
              <th>Precio Unit.</th>
              <th>Cantidad</th>
              <th>Subtotal</th>
            </tr>
          </thead>
          <tbody id="modalSaleProducts">
          </tbody>
        </table>

        <!-- Totals -->
        <div class="row">
          <div class="col-md-4 col-md-offset-8">
            <table class="table">
              <tr>
                <td><strong>Subtotal:</strong></td>
                <td class="text-right" id="modalSaleSubtotal"></td>
              </tr>
              <tr>
                <td><strong>Impuesto:</strong></td>
                <td class="text-right" id="modalSaleTax"></td>
              </tr>
              <tr class="success" style="font-size:18px;">
                <td><strong>TOTAL:</strong></td>
                <td class="text-right" id="modalSaleTotal"></td>
              </tr>
            </table>
          </div>
        </div>

      </div>
      <div class="modal-footer">
        <button type="button" class="btn btn-default" data-dismiss="modal">Cerrar</button>
      </div>
    </div>
  </div>
</div>
