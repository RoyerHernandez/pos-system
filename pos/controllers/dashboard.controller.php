<?php

class DashboardController{

	/*=============================================
	GET TODAY'S SALES
	=============================================*/

	static public function ctrTodaySales(){

		return DashboardModel::mdlTodaySales();

	}

	/*=============================================
	GET LOW STOCK COUNT
	=============================================*/

	static public function ctrLowStockProducts(){

		return DashboardModel::mdlLowStockProducts();

	}

	/*=============================================
	GET LOW STOCK PRODUCTS LIST
	=============================================*/

	static public function ctrLowStockProductsList(){

		return DashboardModel::mdlLowStockProductsList();

	}

	/*=============================================
	GET TOTAL CLIENTS
	=============================================*/

	static public function ctrTotalClients(){

		return DashboardModel::mdlTotalClients();

	}

	/*=============================================
	GET SALES LAST 7 DAYS
	=============================================*/

	static public function ctrSalesLast7Days(){

		return DashboardModel::mdlSalesLast7Days();

	}

	/*=============================================
	GET TOP PRODUCTS TODAY
	=============================================*/

	static public function ctrTopProductsToday(){

		return DashboardModel::mdlTopProductsToday();

	}

	/*=============================================
	GET SALES BY PAYMENT METHOD
	=============================================*/

	static public function ctrSalesByPaymentMethod(){

		return DashboardModel::mdlSalesByPaymentMethod();

	}

}
