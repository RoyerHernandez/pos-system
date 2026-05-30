<?php

session_start();

?>

<!DOCTYPE html>
<html>
<head>

  <meta charset="utf-8">
  <meta http-equiv="X-UA-Compatible" content="IE=edge">

  <title>Inventory System</title>

  <!-- Tell the browser to be responsive to screen width -->
  <meta content="width=device-width, initial-scale=1, maximum-scale=1, user-scalable=no" name="viewport">

  <link rel="icon" href="views/img/plantilla/icono-negro.png">

   <!--=====================================
  CSS PLUGINS
  ======================================-->

  <!-- Bootstrap 3.3.7 -->
  <link rel="stylesheet" href="views/bower_components/bootstrap/dist/css/bootstrap.min.css">

  <!-- Font Awesome -->
  <link rel="stylesheet" href="views/bower_components/font-awesome/css/font-awesome.min.css">

  <!-- Ionicons -->
  <link rel="stylesheet" href="views/bower_components/Ionicons/css/ionicons.min.css">

  <!-- Theme style -->
  <link rel="stylesheet" href="views/dist/css/AdminLTE.css">

  <!-- AdminLTE Skins -->
  <link rel="stylesheet" href="views/dist/css/skins/_all-skins.min.css">

  <!-- Google Font -->
  <link rel="stylesheet" href="https://fonts.googleapis.com/css?family=Source+Sans+Pro:300,400,600,700,300italic,400italic,600italic">

   <!-- DataTables -->
  <link rel="stylesheet" href="views/bower_components/datatables.net-bs/css/dataTables.bootstrap.min.css">
  <link rel="stylesheet" href="views/bower_components/datatables.net-bs/css/responsive.bootstrap.min.css">

  <!--=====================================
  JAVASCRIPT PLUGINS
  ======================================-->

  <!-- jQuery 3 -->
  <script src="views/bower_components/jquery/dist/jquery.min.js"></script>

  <!-- Bootstrap 3.3.7 -->
  <script src="views/bower_components/bootstrap/dist/js/bootstrap.min.js"></script>

  <!-- FastClick -->
  <script src="views/bower_components/fastclick/lib/fastclick.js"></script>

  <!-- AdminLTE App -->
  <script src="views/dist/js/adminlte.min.js"></script>

  <!-- DataTables -->
  <script src="views/bower_components/datatables.net/js/jquery.dataTables.min.js"></script>
  <script src="views/bower_components/datatables.net-bs/js/dataTables.bootstrap.min.js"></script>
  <script src="views/bower_components/datatables.net-bs/js/dataTables.responsive.min.js"></script>
  <script src="views/bower_components/datatables.net-bs/js/responsive.bootstrap.min.js"></script>

  <!-- Chart.js -->
  <script src="views/bower_components/chart.js/Chart.min.js"></script>

  <!-- SweetAlert 2 -->
  <script src="views/plugins/sweetalert2/sweetalert2.all.js"></script>
  <!-- By default SweetAlert2 doesn't support IE. To enable IE 11 support, include Promise polyfill:-->
  <script src="https://cdnjs.cloudflare.com/ajax/libs/core-js/2.4.1/core.js"></script>

</head>

<!--=====================================
DOCUMENT BODY
======================================-->

<body class="hold-transition skin-blue sidebar-collapse sidebar-mini login-page">

  <?php

  if(isset($_SESSION["loggedIn"]) && $_SESSION["loggedIn"] == "ok"){

   echo '<div class="wrapper">';

    /*=============================================
    HEADER
    =============================================*/

    include "modules/header.php";

    /*=============================================
    MENU
    =============================================*/

    include "modules/menu.php";

    /*=============================================
    CONTENT
    =============================================*/

    if(isset($_GET["ruta"])){

      if($_GET["ruta"] == "inicio" ||
         $_GET["ruta"] == "usuarios" ||
         $_GET["ruta"] == "categorias" ||
         $_GET["ruta"] == "productos" ||
         $_GET["ruta"] == "clientes" ||
         $_GET["ruta"] == "ventas" ||
         $_GET["ruta"] == "crear-venta" ||
         $_GET["ruta"] == "caja" ||
         $_GET["ruta"] == "reportes" ||
         $_GET["ruta"] == "salir"){

        include "modules/".$_GET["ruta"].".php";

      }else{

        include "modules/404.php";

      }

    }else{

      include "modules/inicio.php";

    }

    /*=============================================
    FOOTER
    =============================================*/

    include "modules/footer.php";

    echo '</div>';

  }else{

    include "modules/login.php";

  }

  ?>


<script src="views/js/plantilla.js"></script>
<script src="views/js/usuarios.js"></script>
<script src="views/js/sales.js"></script>
<script src="views/js/sales-admin.js"></script>
<script src="views/js/dashboard.js"></script>

</body>
</html>