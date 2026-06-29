<?php

$todaySales = DashboardController::ctrTodaySales();
$lowStock = DashboardController::ctrLowStockProducts();
$lowStockList = DashboardController::ctrLowStockProductsList();
$totalClients = DashboardController::ctrTotalClients();
$cashRegister = CashRegisterModel::mdlGetOpenCashRegister("caja", $_SESSION["id"]);
$salesLast7Days = DashboardController::ctrSalesLast7Days();
$topProducts = DashboardController::ctrTopProductsToday();
$salesByPayment = DashboardController::ctrSalesByPaymentMethod();

// Build 7-day data (fill missing days with 0)
$last7Labels = array();
$last7Data = array();
$salesByDate = array();

foreach($salesLast7Days as $row){
	$salesByDate[$row["sale_date"]] = $row["total_amount"];
}

for($i = 6; $i >= 0; $i--){
	$date = date("Y-m-d", strtotime("-$i days"));
	$label = date("D d", strtotime($date));
	$last7Labels[] = $label;
	$last7Data[] = isset($salesByDate[$date]) ? floatval($salesByDate[$date]) : 0;
}

// Top products data
$topLabels = array();
$topData = array();
foreach($topProducts as $row){
	$topLabels[] = $row["descripcion"];
	$topData[] = intval($row["total_qty"]);
}

// Payment method data
$paymentLabels = array();
$paymentData = array();
$paymentColors = array();
$colorMap = array("Efectivo" => "#00a65a", "Tarjeta" => "#00c0ef", "Transferencia" => "#f39c12", "Mixto" => "#605ca8");

foreach($salesByPayment as $row){
	$paymentLabels[] = $row["metodo_pago"];
	$paymentData[] = floatval($row["total_amount"]);
	$paymentColors[] = isset($colorMap[$row["metodo_pago"]]) ? $colorMap[$row["metodo_pago"]] : "#999";
}

?>

<div class="content-wrapper">

  <section class="content-header">
    <h1>
      Tablero
      <small>Panel de Control</small>
    </h1>
    <ol class="breadcrumb">
      <li><a href="inicio"><i class="fa fa-dashboard"></i> Inicio</a></li>
      <li class="active">Tablero</li>
    </ol>
  </section>

  <section class="content">

    <!--=============================================
    KPI WIDGETS
    =============================================-->

    <div class="row">

      <!-- Today's Sales -->
      <div class="col-lg-3 col-xs-6">
        <div class="small-box bg-green">
          <div class="inner">
            <h3>$<?php echo number_format($todaySales["total_amount"], 2); ?></h3>
            <p>Ventas de Hoy (<?php echo $todaySales["total_count"]; ?>)</p>
          </div>
          <div class="icon">
            <i class="fa fa-shopping-cart"></i>
          </div>
          <a href="ventas" class="small-box-footer">
            Ver ventas <i class="fa fa-arrow-circle-right"></i>
          </a>
        </div>
      </div>

      <!-- Low Stock -->
      <div class="col-lg-3 col-xs-6">
        <div class="small-box bg-red">
          <div class="inner">
            <h3><?php echo $lowStock["total"]; ?></h3>
            <p>Productos Stock Bajo</p>
          </div>
          <div class="icon">
            <i class="fa fa-exclamation-triangle"></i>
          </div>
          <a href="productos" class="small-box-footer">
            Ver productos <i class="fa fa-arrow-circle-right"></i>
          </a>
        </div>
      </div>

      <!-- Cash Register -->
      <div class="col-lg-3 col-xs-6">
        <div class="small-box <?php echo $cashRegister ? 'bg-aqua' : 'bg-yellow'; ?>">
          <div class="inner">
            <?php if($cashRegister): ?>
              <h3>$<?php echo number_format($cashRegister["monto_apertura"] + $cashRegister["total_efectivo"], 2); ?></h3>
              <p>Caja Abierta #<?php echo $cashRegister["id"]; ?></p>
            <?php else: ?>
              <h3>Cerrada</h3>
              <p>Caja</p>
            <?php endif; ?>
          </div>
          <div class="icon">
            <i class="fa fa-desktop"></i>
          </div>
          <a href="caja" class="small-box-footer">
            Gestionar caja <i class="fa fa-arrow-circle-right"></i>
          </a>
        </div>
      </div>

      <!-- Clients -->
      <div class="col-lg-3 col-xs-6">
        <div class="small-box bg-purple">
          <div class="inner">
            <h3><?php echo $totalClients["total"]; ?></h3>
            <p>Clientes Registrados</p>
          </div>
          <div class="icon">
            <i class="fa fa-users"></i>
          </div>
          <a href="clientes" class="small-box-footer">
            Ver clientes <i class="fa fa-arrow-circle-right"></i>
          </a>
        </div>
      </div>

    </div>

    <!--=============================================
    CHARTS ROW 1
    =============================================-->

    <div class="row">

      <!-- Sales Last 7 Days -->
      <div class="col-md-8">
        <div class="box box-primary">
          <div class="box-header with-border">
            <h3 class="box-title"><i class="fa fa-bar-chart"></i> Ventas Ultimos 7 Dias</h3>
          </div>
          <div class="box-body">
            <canvas id="chartSales7Days" height="250"></canvas>
          </div>
        </div>
      </div>

      <!-- Sales by Payment Method -->
      <div class="col-md-4">
        <div class="box box-warning">
          <div class="box-header with-border">
            <h3 class="box-title"><i class="fa fa-pie-chart"></i> Ventas por Metodo de Pago</h3>
            <small class="text-muted pull-right">Hoy</small>
          </div>
          <div class="box-body">
            <?php if(count($paymentData) > 0): ?>
              <canvas id="chartPaymentMethod" height="250"></canvas>
            <?php else: ?>
              <div class="text-center text-muted" style="padding:60px 0;">
                <i class="fa fa-pie-chart fa-3x"></i>
                <p style="margin-top:10px;">Sin ventas hoy</p>
              </div>
            <?php endif; ?>
          </div>
        </div>
      </div>

    </div>

    <!--=============================================
    CHARTS ROW 2
    =============================================-->

    <div class="row">

      <!-- Top 5 Products Today -->
      <div class="col-md-6">
        <div class="box box-success">
          <div class="box-header with-border">
            <h3 class="box-title"><i class="fa fa-trophy"></i> Top 5 Productos del Dia</h3>
          </div>
          <div class="box-body">
            <?php if(count($topData) > 0): ?>
              <canvas id="chartTopProducts" height="220"></canvas>
            <?php else: ?>
              <div class="text-center text-muted" style="padding:50px 0;">
                <i class="fa fa-trophy fa-3x"></i>
                <p style="margin-top:10px;">Sin ventas hoy</p>
              </div>
            <?php endif; ?>
          </div>
        </div>
      </div>

      <!-- Low Stock Alert List -->
      <div class="col-md-6">
        <div class="box box-danger">
          <div class="box-header with-border">
            <h3 class="box-title"><i class="fa fa-exclamation-triangle"></i> Alertas de Stock Bajo</h3>
          </div>
          <div class="box-body no-padding">
            <?php if(count($lowStockList) > 0): ?>
            <table class="table table-striped">
              <thead>
                <tr>
                  <th>Codigo</th>
                  <th>Producto</th>
                  <th>Stock</th>
                  <th>Minimo</th>
                </tr>
              </thead>
              <tbody>
                <?php foreach($lowStockList as $product): ?>
                <tr>
                  <td><small><?php echo $product["codigo"]; ?></small></td>
                  <td><?php echo $product["descripcion"]; ?></td>
                  <td>
                    <span class="label label-<?php echo $product['stock'] == 0 ? 'danger' : 'warning'; ?>">
                      <?php echo $product["stock"]; ?>
                    </span>
                  </td>
                  <td><?php echo $product["stock_minimo"]; ?></td>
                </tr>
                <?php endforeach; ?>
              </tbody>
            </table>
            <?php else: ?>
              <div class="text-center text-muted" style="padding:50px 0;">
                <i class="fa fa-check-circle fa-3x text-success"></i>
                <p style="margin-top:10px;">Todos los productos tienen stock suficiente</p>
              </div>
            <?php endif; ?>
          </div>
        </div>
      </div>

    </div>

  </section>

</div>

<!--=============================================
CHART.JS DATA
=============================================-->

<script>

var dashboardData = {
	sales7Days: {
		labels: <?php echo json_encode($last7Labels); ?>,
		data: <?php echo json_encode($last7Data); ?>
	},
	topProducts: {
		labels: <?php echo json_encode($topLabels); ?>,
		data: <?php echo json_encode($topData); ?>
	},
	paymentMethod: {
		labels: <?php echo json_encode($paymentLabels); ?>,
		data: <?php echo json_encode($paymentData); ?>,
		colors: <?php echo json_encode($paymentColors); ?>
	}
};

</script>
