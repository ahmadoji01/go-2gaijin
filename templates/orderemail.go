package templates

func DeliveryEmail(itemName string, name string, email string, phone string, wechat string, facebook string, origin string, destination string, deliveryTime string, notes string) string {
	template := `<!DOCTYPE html>
	<html>
	<head>
	  <meta charset="utf-8">
	  <meta http-equiv="x-ua-compatible" content="ie=edge">
	  <title>Email Receipt</title>
	  <meta name="viewport" content="width=device-width, initial-scale=1">
	  <style type="text/css">
	  @media screen {
		@font-face {
		  font-family: 'Source Sans Pro';
		  font-style: normal;
		  font-weight: 400;
		  src: local('Source Sans Pro Regular'), local('SourceSansPro-Regular'), url(https://fonts.gstatic.com/s/sourcesanspro/v10/ODelI1aHBYDBqgeIAH2zlBM0YzuT7MdOe03otPbuUS0.woff) format('woff');
		}
	
		@font-face {
		  font-family: 'Source Sans Pro';
		  font-style: normal;
		  font-weight: 700;
		  src: local('Source Sans Pro Bold'), local('SourceSansPro-Bold'), url(https://fonts.gstatic.com/s/sourcesanspro/v10/toadOcfmlt9b38dHJxOBGFkQc6VGVFSmCnC_l7QZG60.woff) format('woff');
		}
	  }
	  body,
	  table,
	  td,
	  a {
		-ms-text-size-adjust: 100%; /* 1 */
		-webkit-text-size-adjust: 100%; /* 2 */
	  }
	  table,
	  td {
		mso-table-rspace: 0pt;
		mso-table-lspace: 0pt;
	  }
	  img {
		-ms-interpolation-mode: bicubic;
	  }
	  a[x-apple-data-detectors] {
		font-family: inherit !important;
		font-size: inherit !important;
		font-weight: inherit !important;
		line-height: inherit !important;
		color: inherit !important;
		text-decoration: none !important;
	  }
	  div[style*="margin: 16px 0;"] {
		margin: 0 !important;
	  }
	
	  body {
		width: 100% !important;
		height: 100% !important;
		padding: 0 !important;
		margin: 0 !important;
	  }
	  table {
		border-collapse: collapse !important;
	  }
	
	  a {
		color: #1a82e2;
	  }
	
	  img {
		height: auto;
		line-height: 100%;
		text-decoration: none;
		border: 0;
		outline: none;
	  }
	  </style>
	
	</head>
	<body style="background-color: #D2C7BA;">
	
	  <!-- start preheader -->
	  <div class="preheader" style="display: none; max-width: 0; max-height: 0; overflow: hidden; font-size: 1px; line-height: 1px; color: #fff; opacity: 0;">
		Below are the request details
	  </div>
	  <!-- end preheader -->
	
	  <!-- start body -->
	  <table border="0" cellpadding="0" cellspacing="0" width="100%">
	
		<!-- start logo -->
		<tr>
		  <td align="center" bgcolor="#D2C7BA">
			<!--[if (gte mso 9)|(IE)]>
			<table align="center" border="0" cellpadding="0" cellspacing="0" width="600">
			<tr>
			<td align="center" valign="top" width="600">
			<![endif]-->
			<table border="0" cellpadding="0" cellspacing="0" width="100%" style="max-width: 600px;">
			  <tr>
				<td align="center" valign="top" style="padding: 36px 24px;">
				  <a href="https://2gaijin.com" target="_blank" style="display: inline-block;">
					<img src="https://storage.googleapis.com/rails-2gaijin-storage/2gaijinheader.png" alt="Logo" border="0" width="48" style="display: block; width: 256px; max-width: 256px; min-width: 256px;">
				  </a>
				</td>
			  </tr>
			</table>
			<!--[if (gte mso 9)|(IE)]>
			</td>
			</tr>
			</table>
			<![endif]-->
		  </td>
		</tr>
		<!-- end logo -->
	
		<!-- start hero -->
		<tr>
		  <td align="center" bgcolor="#D2C7BA">
			<!--[if (gte mso 9)|(IE)]>
			<table align="center" border="0" cellpadding="0" cellspacing="0" width="600">
			<tr>
			<td align="center" valign="top" width="600">
			<![endif]-->
			<table border="0" cellpadding="0" cellspacing="0" width="100%" style="max-width: 600px;">
			  <tr>
				<td align="left" bgcolor="#ffffff" style="padding: 36px 24px 0; font-family: 'Source Sans Pro', Helvetica, Arial, sans-serif; border-top: 3px solid #d4dadf;">
				  <h1 style="margin: 0; font-size: 32px; font-weight: 700; letter-spacing: -1px; line-height: 48px;">Thank you for your order!</h1>
				</td>
			  </tr>
			</table>
			<!--[if (gte mso 9)|(IE)]>
			</td>
			</tr>
			</table>
			<![endif]-->
		  </td>
		</tr>
		<!-- end hero -->
	
		<!-- start copy block -->
		<tr>
		  <td align="center" bgcolor="#D2C7BA">
			<!--[if (gte mso 9)|(IE)]>
			<table align="center" border="0" cellpadding="0" cellspacing="0" width="600">
			<tr>
			<td align="center" valign="top" width="600">
			<![endif]-->
			<table border="0" cellpadding="0" cellspacing="0" width="100%" style="max-width: 600px;">
	
			  <!-- start copy -->
			  <tr>
				<td align="left" bgcolor="#ffffff" style="padding: 24px; font-family: 'Source Sans Pro', Helvetica, Arial, sans-serif; font-size: 16px; line-height: 24px;">
				  <p style="margin: 0;">Here is a summary of the delivery request</p>
				</td>
			  </tr>
			  <!-- end copy -->

			  	<!-- start customer's info table -->
				<tr>
					<td align="left" bgcolor="#ffffff" style="padding: 24px; font-family: 'Source Sans Pro', Helvetica, Arial, sans-serif; font-size: 16px; line-height: 24px;">
					<table border="0" cellpadding="0" cellspacing="0" width="100%">
						<tr>
						<td align="left" bgcolor="#D2C7BA" width="25%" style="padding: 12px;font-family: 'Source Sans Pro', Helvetica, Arial, sans-serif; font-size: 16px; line-height: 24px;"><strong>Customer's Info</strong></td>
						<td align="left" bgcolor="#D2C7BA" width="75%" style="padding: 12px;font-family: 'Source Sans Pro', Helvetica, Arial, sans-serif; font-size: 16px; line-height: 24px;"></td>
						</tr>
						<tr>
						<td align="left" width="25%" style="padding: 6px 12px;font-family: 'Source Sans Pro', Helvetica, Arial, sans-serif; font-size: 16px; line-height: 24px;">Name</td>
						<td align="left" width="75%" style="padding: 6px 12px;font-family: 'Source Sans Pro', Helvetica, Arial, sans-serif; font-size: 16px; line-height: 24px;">` + name + `</td>
						</tr>
						<tr>
						<td align="left" width="25%" style="padding: 6px 12px;font-family: 'Source Sans Pro', Helvetica, Arial, sans-serif; font-size: 16px; line-height: 24px;">Email</td>
						<td align="left" width="75%" style="padding: 6px 12px;font-family: 'Source Sans Pro', Helvetica, Arial, sans-serif; font-size: 16px; line-height: 24px;">` + email + `</td>
						</tr>
						<tr>
						<td align="left" width="25%" style="padding: 6px 12px;font-family: 'Source Sans Pro', Helvetica, Arial, sans-serif; font-size: 16px; line-height: 24px;">Phone</td>
						<td align="left" width="75%" style="padding: 6px 12px;font-family: 'Source Sans Pro', Helvetica, Arial, sans-serif; font-size: 16px; line-height: 24px;">` + phone + `</td>
						</tr>
						<tr>
						<td align="left" width="25%" style="padding: 6px 12px;font-family: 'Source Sans Pro', Helvetica, Arial, sans-serif; font-size: 16px; line-height: 24px;">Facebook</td>
						<td align="left" width="75%" style="padding: 6px 12px;font-family: 'Source Sans Pro', Helvetica, Arial, sans-serif; font-size: 16px; line-height: 24px;">` + facebook + `</td>
						</tr>
						<tr>
						<td align="left" width="25%" style="padding: 6px 12px;font-family: 'Source Sans Pro', Helvetica, Arial, sans-serif; font-size: 16px; line-height: 24px;">WeChat ID</td>
						<td align="left" width="75%" style="padding: 6px 12px;font-family: 'Source Sans Pro', Helvetica, Arial, sans-serif; font-size: 16px; line-height: 24px;">` + wechat + `</td>
						</tr>
						<tr>
						<td align="left" width="25%" style="padding: 12px; font-family: 'Source Sans Pro', Helvetica, Arial, sans-serif; font-size: 16px; line-height: 24px; border-top: 2px dashed #D2C7BA; border-bottom: 2px dashed #D2C7BA;"></td>
						<td align="left" width="75%" style="padding: 12px; font-family: 'Source Sans Pro', Helvetica, Arial, sans-serif; font-size: 16px; line-height: 24px; border-top: 2px dashed #D2C7BA; border-bottom: 2px dashed #D2C7BA;"></td>
						</tr>
					</table>
					</td>
				</tr>
				<!-- end customer's info table -->
	
			  <!-- start receipt table -->
			  <tr>
				<td align="left" bgcolor="#ffffff" style="padding: 24px; font-family: 'Source Sans Pro', Helvetica, Arial, sans-serif; font-size: 16px; line-height: 24px;">
				  <table border="0" cellpadding="0" cellspacing="0" width="100%">
					<tr>
					  <td align="left" bgcolor="#D2C7BA" width="70%" style="padding: 12px;font-family: 'Source Sans Pro', Helvetica, Arial, sans-serif; font-size: 16px; line-height: 24px;"><strong>Item Name</strong></td>
					  <td align="left" bgcolor="#D2C7BA" width="30%" style="padding: 12px;font-family: 'Source Sans Pro', Helvetica, Arial, sans-serif; font-size: 16px; line-height: 24px;"><strong></strong></td>
					</tr>
					<tr>
						<td align="left" width="70%" style="padding: 6px 12px;font-family: 'Source Sans Pro', Helvetica, Arial, sans-serif; font-size: 16px; line-height: 24px;">` + itemName + `</td>
						<td align="left" width="30%" style="padding: 6px 12px;font-family: 'Source Sans Pro', Helvetica, Arial, sans-serif; font-size: 16px; line-height: 24px;"></td>
					</tr>
					<tr>
					  <td align="left" width="70%" style="padding: 12px; font-family: 'Source Sans Pro', Helvetica, Arial, sans-serif; font-size: 16px; line-height: 24px; border-top: 2px dashed #D2C7BA; border-bottom: 2px dashed #D2C7BA;"></td>
					  <td align="left" width="30%" style="padding: 12px; font-family: 'Source Sans Pro', Helvetica, Arial, sans-serif; font-size: 16px; line-height: 24px; border-top: 2px dashed #D2C7BA; border-bottom: 2px dashed #D2C7BA;"></td>
					</tr>
				  </table>
				</td>
			  </tr>
			  <!-- end reeipt table -->

			  <!-- start notes table -->
			  <tr>
				<td align="left" bgcolor="#ffffff" style="padding: 24px; font-family: 'Source Sans Pro', Helvetica, Arial, sans-serif; font-size: 16px; line-height: 24px;">
				  <table border="0" cellpadding="0" cellspacing="0" width="100%">
					<tr>
					  <td align="left" bgcolor="#D2C7BA" width="70%" style="padding: 12px;font-family: 'Source Sans Pro', Helvetica, Arial, sans-serif; font-size: 16px; line-height: 24px;"><strong>Delivery Notes</strong></td>
					  <td align="left" bgcolor="#D2C7BA" width="30%" style="padding: 12px;font-family: 'Source Sans Pro', Helvetica, Arial, sans-serif; font-size: 16px; line-height: 24px;"><strong></strong></td>
					</tr>
					<tr>
						<td align="left" width="70%" style="padding: 6px 12px;font-family: 'Source Sans Pro', Helvetica, Arial, sans-serif; font-size: 16px; line-height: 24px;">` + notes + `</td>
						<td align="left" width="30%" style="padding: 6px 12px;font-family: 'Source Sans Pro', Helvetica, Arial, sans-serif; font-size: 16px; line-height: 24px;"></td>
					</tr>
					<tr>
					  <td align="left" width="70%" style="padding: 12px; font-family: 'Source Sans Pro', Helvetica, Arial, sans-serif; font-size: 16px; line-height: 24px; border-top: 2px dashed #D2C7BA; border-bottom: 2px dashed #D2C7BA;"></td>
					  <td align="left" width="30%" style="padding: 12px; font-family: 'Source Sans Pro', Helvetica, Arial, sans-serif; font-size: 16px; line-height: 24px; border-top: 2px dashed #D2C7BA; border-bottom: 2px dashed #D2C7BA;"></td>
					</tr>
				  </table>
				</td>
			  </tr>
			  <!-- end notes table -->
	
			</table>
			<!--[if (gte mso 9)|(IE)]>
			</td>
			</tr>
			</table>
			<![endif]-->
		  </td>
		</tr>
		<!-- end copy block -->
	
		<!-- start receipt address block -->
		<tr>
		  <td align="center" bgcolor="#D2C7BA" valign="top" width="100%">
			<!--[if (gte mso 9)|(IE)]>
			<table align="center" border="0" cellpadding="0" cellspacing="0" width="600">
			<tr>
			<td align="center" valign="top" width="600">
			<![endif]-->
			<table align="center" bgcolor="#ffffff" border="0" cellpadding="0" cellspacing="0" width="100%" style="max-width: 600px;">
			  <tr>
				<td align="center" valign="top" style="font-size: 0; border-bottom: 3px solid #d4dadf">
				  <!--[if (gte mso 9)|(IE)]>
				  <table align="center" border="0" cellpadding="0" cellspacing="0" width="600">
				  <tr>
				  <td align="left" valign="top" width="300">
				  <![endif]-->
				  <div style="display: inline-block; width: 100%; max-width: 50%; min-width: 240px; vertical-align: top;">
					<table align="left" border="0" cellpadding="0" cellspacing="0" width="100%" style="max-width: 300px;">
					  <tr>
						<td align="left" valign="top" style="padding-bottom: 36px; padding-left: 36px; font-family: 'Source Sans Pro', Helvetica, Arial, sans-serif; font-size: 16px; line-height: 24px;">
						  <p><strong>Destination</strong></p>
						  <p>` + destination + `</p>
						</td>
					  </tr>
					  <tr>
						<td align="left" valign="top" style="padding-bottom: 36px; padding-left: 36px; font-family: 'Source Sans Pro', Helvetica, Arial, sans-serif; font-size: 16px; line-height: 24px;">
						  <p><strong>Origin</strong></p>
						  <p>` + origin + `</p>
						</td>
					  </tr>
					</table>
				  </div>
				  <div style="display: inline-block; width: 100%; max-width: 50%; min-width: 240px; vertical-align: top;">
					<table align="left" border="0" cellpadding="0" cellspacing="0" width="100%" style="max-width: 300px;">
					  <tr>
						<td align="left" valign="top" style="padding-bottom: 36px; padding-left: 36px; font-family: 'Source Sans Pro', Helvetica, Arial, sans-serif; font-size: 16px; line-height: 24px;">
						  <p><strong>Delivery Date</strong></p>
						  <p>` + deliveryTime + `</p>
						</td>
					  </tr>
					</table>
				  </div>
				</td>
			  </tr>
			</table>
			<!--[if (gte mso 9)|(IE)]>
			</td>
			</tr>
			</table>
			<![endif]-->
		  </td>
		</tr>
		<!-- end receipt address block -->
	
		<!-- start footer -->
		<tr>
		  <td align="center" bgcolor="#D2C7BA" style="padding: 24px;">
			<!--[if (gte mso 9)|(IE)]>
			<table align="center" border="0" cellpadding="0" cellspacing="0" width="600">
			<tr>
			<td align="center" valign="top" width="600">
			<![endif]-->
			<table border="0" cellpadding="0" cellspacing="0" width="100%" style="max-width: 600px;">
	
			  <!-- start unsubscribe -->
			  <tr>
				<td align="center" bgcolor="#D2C7BA" style="padding: 12px 24px; font-family: 'Source Sans Pro', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 20px; color: #666;">
				  <h4>2Gaijin.com</h4>
				  <p style="margin: 0;">Kita-ku, Sapporo, Hokkaido, Japan</p>
				</td>
			  </tr>
			  <!-- end unsubscribe -->
	
			</table>
			<!--[if (gte mso 9)|(IE)]>
			</td>
			</tr>
			</table>
			<![endif]-->
		  </td>
		</tr>
		<!-- end footer -->
	
	  </table>
	  <!-- end body -->
	
	</body>
	</html>
	`
	return template
}

func OrderNotificationEmail(link string, itemName string) string {
	template := `<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01 Transitional//EN" "http://www.w3.org/TR/html4/loose.dtd">
	<html xmlns="http://www.w3.org/1999/xhtml" xmlns:v="urn:schemas-microsoft-com:vml" xmlns:o="urn:schemas-microsoft-com:office:office">
	  <head>
		<meta content="text/html; charset=UTF-8" http-equiv="Content-Type" />
		<!-- [ if !mso]> <!-->
		<meta content="IE=edge" http-equiv="X-UA-Compatible" />
		<!-- <![endif] -->
		<meta content="telephone=no" name="format-detection" />
		<meta content="width=device-width, initial-scale=1.0" name="viewport" />
		<link rel="apple-touch-icon" sizes="76x76" href="https://storage.cloud.google.com/rails-2gaijin-storage/2gaijinicon.png">
		<link rel="icon" type="image/png" sizes="96x96" href="https://storage.cloud.google.com/rails-2gaijin-storage/2gaijinicon.png">
		<title>2Gaijin.com: We are migrating to new website!</title>
		<link href="https://fonts.googleapis.com/css?family=Roboto+Mono" rel="stylesheet">
		<script src='http://cdnjs.cloudflare.com/ajax/libs/jquery/2.1.3/jquery.min.js'></script>
		<script src="http://paulgoddarddesign.com/js/ripple.js"></script>
		<style type="text/css">
		  .ExternalClass {width: 100%;}
		  .ExternalClass, .ExternalClass p, .ExternalClass span, .ExternalClass font, .ExternalClass td, .ExternalClass div, .ExternalClass b, .ExternalClass br, .ExternalClass img {line-height: 100% !important;}
		  /* iOS BLUE LINKS */
		  .appleBody a {color:#212121; text-decoration: none;}
		  .appleFooter a {color:#212121!important; text-decoration: none!important;}
		  /* END iOS BLUE LINKS */
		  img {color: #ffffff;text-align: center;font-family: Open Sans, Helvetica, Arial, sans-serif;display: block;}
		  body {margin: 0;padding: 0;-webkit-text-size-adjust: 100% !important;-ms-text-size-adjust: 100% !important;font-family: 'Open Sans', Helvetica, Arial, sans-serif!important;}
		  body,#body_style {background: #fffffe;}
		  table td {border-collapse: collapse;border-spacing: 0 !important;}
		  table tr {border-collapse: collapse;border-spacing: 0 !important;}
		  table tbody {border-collapse: collapse;border-spacing: 0 !important;}
		  table {border-collapse: collapse;border-spacing: 0 !important;}
		  span.yshortcuts,a span.yshortcuts {color: #000001;background-color: none;border: none;}
		  span.yshortcuts:hover,
		  span.yshortcuts:active,
		  span.yshortcuts:focus {color: #000001; background-color: none; border: none;}
		  img {-ms-interpolation-mode: : bicubic;}
		  a[x-apple-data-detectors] {color: inherit !important;text-decoration: none !important;font-size: inherit !important;font-family: inherit !important;font-weight: inherit !important;line-height: inherit !important;
		  }
		  /**** My desktop styles ****/
		  @media only screen and (min-width: 600px) {
			.noDesk {display: none !important;}
			.td-padding {padding-left: 15px!important;padding-right: 15px!important;}
			.padding-container {padding: 0px 15px 0px 15px!important;mso-padding-alt: 0px 15px 0px 15px!important;}
			.mobile-column-left-padding { padding: 0px 0px 0px 0px!important; mso-alt-padding: 0px 0px 0px 0px!important; }
			.mobile-column-right-padding { padding: 0px 0px 0px 0px!important; mso-alt-padding: 0px 0px 0px 0px!important; }
			.mobile {display: none !important}
		  }
		  /**** My mobile styles ****/
		  @media only screen and (max-width: 599px) and (-webkit-min-device-pixel-ratio: 1) {
			*[class].wrapper { width:100% !important; }
			*[class].container { width:100% !important; }
			*[class].mobile { width:100% !important; display:block !important; }
			*[class].image{ width:100% !important; height:auto; }
			*[class].center{ margin:0 auto !important; text-align:center !important; }
			*[class="mobileOff"] { width: 0px !important; display: none !important; }
			*[class*="mobileOn"] { display: block !important; max-height:none !important; }
			p[class="mobile-padding"] {padding-left: 0px!important;padding-top: 10px;}
			.padding-container {padding: 0px 15px 0px 15px!important;mso-padding-alt: 0px 15px 0px 15px!important;}
			.hund {width: 100% !important;height: auto !important;}
			.td-padding {padding-left: 15px!important;padding-right: 15px!important;}
			.mobile-column-left-padding { padding: 18px 0px 18px 0px!important; mso-alt-padding: 18px 0px 18px 0px!important; }
			.mobile-column-right-padding { padding: 18px 0px 0px 0px!important; mso-alt-padding: 18px 0px 0px 0px!important; }
			.stack { width: 100% !important; }
			img {width: 100%!important;height: auto!important;}
			*[class="hide"] {display: none !important}
			*[class="Gmail"] {display: none !important}
			.Gmail {display: none !important}
			.bottom-padding-fix {padding: 0px 0px 18px 0px!important; mso-alt-padding: 0px 0px 18px 0px;}
		  }
		  .social, .social:active {
			opacity: 1!important;
			transform: scale(1);
			transition: all .2s!important;
		  }
		  .social:hover {
			opacity: 0.8!important;
			transform: scale(1.1);
			transition: all .2s!important;
		  }
		  .button.raised {
			transition: box-shadow 0.2s cubic-bezier(0.4, 0, 0.2, 1);
			transition: all .2s;box-shadow: 0 2px 5px 0 rgba(0, 0, 0, 0.26);
		  }
		  .button.raised:hover {
			box-shadow: 0 8px 17px 0 rgba(0, 0, 0, 0.2);transition: all .2s;
			-webkit-box-shadow: 0 8px 17px 0 rgba(0, 0, 0, 0.2);transition: all .2s;
			-moz-box-shadow: 0 8px 17px 0 rgba(0, 0, 0, 0.2);transition: all .2s;
		  }
		  .card-1 {
			box-shadow: 0 2px 5px 0 rgba(0,0,0,.16), 0 2px 10px 0 rgba(0,0,0,.12);
			-webkit-box-shadow: 0 2px 5px 0 rgba(0,0,0,.16), 0 2px 10px 0 rgba(0,0,0,.12);
			-moz-box-shadow: 0 2px 5px 0 rgba(0,0,0,.16), 0 2px 10px 0 rgba(0,0,0,.12);
			transition: box-shadow .45s;
		  }
		  .card-1:hover {
			box-shadow: 0 8px 17px 0 rgba(0,0,0,.2), 0 6px 20px 0 rgba(0,0,0,.19);
			-webkit-box-shadow: 0 8px 17px 0 rgba(0,0,0,.2), 0 6px 20px 0 rgba(0,0,0,.19);
			-moz-box-shadow: 0 8px 17px 0 rgba(0,0,0,.2), 0 6px 20px 0 rgba(0,0,0,.19);
			transition: box-shadow .45s;
		  }
		  .ripplelink{
			display:block
			color:#fff;
			text-decoration:none;
			position:relative;
			overflow:hidden;
			-webkit-transition: all 0.2s ease;
			-moz-transition: all 0.2s ease;
			-o-transition: all 0.2s ease;
			transition: all 0.2s ease;
			z-index:0;
		  }
		  .ripplelink:hover{
			z-index:1000;
		  }
		  .ink {
			display: block;
			position: absolute;
			background:rgba(255, 255, 255, 0.3);
			border-radius: 100%;
			-webkit-transform:scale(0);
			   -moz-transform:scale(0);
				 -o-transform:scale(0);
					transform:scale(0);
		  }
		  .animate {
			-webkit-animation:ripple 0.65s linear;
			 -moz-animation:ripple 0.65s linear;
			  -ms-animation:ripple 0.65s linear;
			   -o-animation:ripple 0.65s linear;
				  animation:ripple 0.65s linear;
		  }
		  @-webkit-keyframes ripple {
			  100% {opacity: 0; -webkit-transform: scale(2.5);}
		  }
		  @-moz-keyframes ripple {
			  100% {opacity: 0; -moz-transform: scale(2.5);}
		  }
		  @-o-keyframes ripple {
			  100% {opacity: 0; -o-transform: scale(2.5);}
		  }
		  @keyframes ripple {
			  100% {opacity: 0; transform: scale(2.5);}
		  }
	
		</style>
		<!--[if gte mso 9]>
		<xml>
		  <o:OfficeDocumentSettings>
			<o:AllowPNG/>
			<o:PixelsPerInch>96</o:PixelsPerInch>
		  </o:OfficeDocumentSettings>
		</xml>
		<![endif]-->
	  </head>
	  <body style="margin:0; padding:0; background-color: #eeeeee;" bgcolor="#eeeeee">
		<!--[if mso]>
		<style type="text/css">
		body, table, td {font-family: Arial, Helvetica, sans-serif !important;}
		</style>
		<![endif]-->
		<!-- START EMAIL -->
		<table width="100%" cellpadding="0" cellspacing="0" border="0" bgcolor="#eeeeee">
		  <div class="Gmail" style="height: 1px !important; margin-top: -1px !important; max-width: 600px !important; min-width: 600px !important; width: 600px !important;"></div>
		  <div style="display: none; max-height: 0px; overflow: hidden;">
			
		  </div>
		  <!-- Insert &zwnj;&nbsp; hack after hidden preview text -->
		  <div style="display: none; max-height: 0px; overflow: hidden;">
			&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
		  </div>
	
		  <!-- START LOGO -->
		  <tr>
			<td width="100%" valign="top" align="center" class="padding-container" style="padding: 18px 0px 18px 0px!important; mso-padding-alt: 18px 0px 18px 0px;">
			  <table width="600" cellpadding="0" cellspacing="0" border="0" align="center" class="wrapper">
				<tr>
				  <td align="center">
					<table cellpadding="0" cellspacing="0" border="0">
					  <tr>
						<td width="100%" valign="top" align="center">
						  <table width="600" cellpadding="0" cellspacing="0" border="0" align="center" class="wrapper" bgcolor="#eeeeee">
							<tr>
							  <td align="center">
								<table width="600" cellpadding="0" cellspacing="0" border="0" class="container" align="center">
								  <!-- START HEADER IMAGE -->
								  <tr>
									<td align="center" class="hund" width="600">
									  <img src="https://storage.googleapis.com/rails-2gaijin-storage/2gaijinheader.png" width="300" alt="Logo" border="0" style="max-width: 300px; display:block;">
									</td>
								  </tr>
								  <!-- END HEADER IMAGE -->
								</table>
							  </td>
							</tr>
						  </table>
						</td>
					  </tr>
					</table>
				  </td>
				</tr>
			  </table>
			</td>
		  </tr>
		  <!-- END LOGO -->
	
		  <!-- START CARD 1 -->
		  <tr>
			<td width="100%" valign="top" align="center" class="padding-container" style="padding-top: 0px!important; padding-bottom: 18px!important; mso-padding-alt: 0px 0px 18px 0px;">
			  <table width="600" cellpadding="0" cellspacing="0" border="0" align="center" class="wrapper">
				<tr>
				  <td>
					<table cellpadding="0" cellspacing="0" border="0">
					  <tr>
						<td style="border-radius: 3px; border-bottom: 2px solid #d4d4d4;" class="card-1" width="100%" valign="top" align="center">
						  <table style="border-radius: 3px;" width="600" cellpadding="0" cellspacing="0" border="0" align="center" class="wrapper" bgcolor="#ffffff">
							<tr>
							  <td align="center">
								<table width="600" cellpadding="0" cellspacing="0" border="0" class="container">
								  <!-- START HEADER IMAGE -->
								  <tr>
									<td align="center" class="hund ripplelink" width="400" style="padding-top: 20px;">
									  <img align="center" width="200" style="border-radius: 3px 3px 0px 0px; width: 100%; max-width: 200px!important" class="hund" src="https://img.icons8.com/color/452/talk-male--v2.png">
									</td>
								  </tr>
								  <!-- END HEADER IMAGE -->
								  <!-- START BODY COPY -->
								  <tr>
									<td class="td-padding" align="left" style="font-family: 'Roboto Mono', monospace; color: #212121!important; font-size: 24px; line-height: 30px; padding-top: 18px; padding-left: 18px!important; padding-right: 18px!important; padding-bottom: 0px!important; mso-line-height-rule: exactly; mso-padding-alt: 18px 18px 0px 13px;">
									  Someone just requested your ` + itemName + `
									</td>
								  </tr>
								  <!-- END BODY COPY -->
								  <!-- BUTTON -->
								  <tr>
									<td align="center" style="padding: 18px 18px 18px 18px; mso-alt-padding: 18px 18px 18px 18px!important;">
									  <a class="button raised" href="` + link + `" target="_blank" style="font-size: 24px; line-height: 14px; font-weight: 500; font-family: Helvetica, Arial, sans-serif; background-color: #ff8c00; color: #ffffff; text-decoration: none; border-radius: 3px; padding: 25px 100px; border: 1px solid #17bef7; display: inline-block;">Click Here to Review the Request</a>
									  <!-- <table width="100%" border="0" cellspacing="0" cellpadding="0">
										<tr>
										  <td>
											<table border="0" cellspacing="0" cellpadding="0">
											  <tr>
												<td align="center" style="border-radius: 3px;" bgcolor="#17bef7">
												  
												</td>
											  </tr>
											</table>
										  </td>
										</tr>
									  </table> -->
									</td>
								  </tr>
								  <!-- END BUTTON -->
								</table>
							  </td>
							</tr>
						  </table>
						</td>
					  </tr>
					</table>
				  </td>
				</tr>
			  </table>
			</td>
		  </tr>
		  <!-- END CARD 1 -->
		  <!-- FOOTER -->
		  <tr>
			<td width="100%" valign="top" align="center" class="padding-container">
			  <table width="600" cellpadding="0" cellspacing="0" border="0" align="center" class="wrapper">
				<tr>
				  <td width="100%" valign="top" align="center">
					<table width="600" cellpadding="0" cellspacing="0" border="0" align="center" class="wrapper" bgcolor="#eeeeee">
					  <tr>
						<td align="center">
						  <table width="600" cellpadding="0" cellspacing="0" border="0" class="container">
							<tr>
							  <!-- SOCIAL -->
							  <td align="center" width="300" style="padding-top: 0px!important; padding-bottom: 18px!important; mso-padding-alt: 0px 0px 18px 0px;">
								<table border="0" cellspacing="0" cellpadding="0">
								  <tr>
									<td align="right" valign="top" class="social">
									  <a href="https://www.facebook.com/2gaijin/"
									  target="_blank">
									  <img src="http://paulgoddarddesign.com/emails/images/material-design/fb-icon.png"
									  height="24" alt="Facebook" border="0" style="display:block; max-width: 24px">
									  </a>
									</td>
									<td width="20"></td>
									<td align="right" valign="top" class="social">
									  <a href="mailto:2gaijin@kitalabs.com"
									  target="_blank">
									  <img src="https://cdn4.iconfinder.com/data/icons/new-google-logo-2015/400/new-google-favicon-512.png"
									  height="24" alt="Google" border="0" style="display:block; max-width: 24px">
									  </a>
									</td>
								  </tr>
								</table>
							  </td>
							  <!-- END SOCIAL -->
							</tr>
							<tr>
							  <td class="td-padding" align="center" style="font-family: 'Roboto Mono', monospace; color: #212121!important; font-size: 16px; line-height: 24px; padding-top: 0px; padding-left: 0px!important; padding-right: 0px!important; padding-bottom: 0px!important; mso-line-height-rule: exactly; mso-padding-alt: 0px 0px 0px 0px;">
								&copy; 2019 2Gaijin.com
							  </td>
							</tr>
						  </table>
						</td>
					  </tr>
					</table>
				  </td>
				</tr>
			  </table>
			</td>
		  </tr>
		  <!-- FOOTER -->
	
		  <!-- SPACER -->
		  <!--[if gte mso 9]>
		  <table align="center" border="0" cellspacing="0" cellpadding="0" width="600">
			<tr>
			  <td align="center" valign="top" width="600" height="36">
				<![endif]-->
				<tr><td height="36px"></td></tr>
				<!--[if gte mso 9]>
			  </td>
			</tr>
		  </table>
		  <![endif]-->
		  <!-- END SPACER -->
	
		</table>
		<!-- END EMAIL -->
		<div style="display:none; white-space:nowrap; font:15px courier; line-height:0;">
		  &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp;
		  &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp;
		  &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp;
		</div>
	  </body>
	</html>
	`
	return template
}

func OrderAcceptedEmail(itemName string) string {
	template := `<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01 Transitional//EN" "http://www.w3.org/TR/html4/loose.dtd">
	<html xmlns="http://www.w3.org/1999/xhtml" xmlns:v="urn:schemas-microsoft-com:vml" xmlns:o="urn:schemas-microsoft-com:office:office">
	  <head>
		<meta content="text/html; charset=UTF-8" http-equiv="Content-Type" />
		<!-- [ if !mso]> <!-->
		<meta content="IE=edge" http-equiv="X-UA-Compatible" />
		<!-- <![endif] -->
		<meta content="telephone=no" name="format-detection" />
		<meta content="width=device-width, initial-scale=1.0" name="viewport" />
		<link rel="apple-touch-icon" sizes="76x76" href="https://storage.cloud.google.com/rails-2gaijin-storage/2gaijinicon.png">
		<link rel="icon" type="image/png" sizes="96x96" href="https://storage.cloud.google.com/rails-2gaijin-storage/2gaijinicon.png">
		<title>2Gaijin.com: We are migrating to new website!</title>
		<link href="https://fonts.googleapis.com/css?family=Roboto+Mono" rel="stylesheet">
		<script src='http://cdnjs.cloudflare.com/ajax/libs/jquery/2.1.3/jquery.min.js'></script>
		<script src="http://paulgoddarddesign.com/js/ripple.js"></script>
		<style type="text/css">
		  .ExternalClass {width: 100%;}
		  .ExternalClass, .ExternalClass p, .ExternalClass span, .ExternalClass font, .ExternalClass td, .ExternalClass div, .ExternalClass b, .ExternalClass br, .ExternalClass img {line-height: 100% !important;}
		  /* iOS BLUE LINKS */
		  .appleBody a {color:#212121; text-decoration: none;}
		  .appleFooter a {color:#212121!important; text-decoration: none!important;}
		  /* END iOS BLUE LINKS */
		  img {color: #ffffff;text-align: center;font-family: Open Sans, Helvetica, Arial, sans-serif;display: block;}
		  body {margin: 0;padding: 0;-webkit-text-size-adjust: 100% !important;-ms-text-size-adjust: 100% !important;font-family: 'Open Sans', Helvetica, Arial, sans-serif!important;}
		  body,#body_style {background: #fffffe;}
		  table td {border-collapse: collapse;border-spacing: 0 !important;}
		  table tr {border-collapse: collapse;border-spacing: 0 !important;}
		  table tbody {border-collapse: collapse;border-spacing: 0 !important;}
		  table {border-collapse: collapse;border-spacing: 0 !important;}
		  span.yshortcuts,a span.yshortcuts {color: #000001;background-color: none;border: none;}
		  span.yshortcuts:hover,
		  span.yshortcuts:active,
		  span.yshortcuts:focus {color: #000001; background-color: none; border: none;}
		  img {-ms-interpolation-mode: : bicubic;}
		  a[x-apple-data-detectors] {color: inherit !important;text-decoration: none !important;font-size: inherit !important;font-family: inherit !important;font-weight: inherit !important;line-height: inherit !important;
		  }
		  /**** My desktop styles ****/
		  @media only screen and (min-width: 600px) {
			.noDesk {display: none !important;}
			.td-padding {padding-left: 15px!important;padding-right: 15px!important;}
			.padding-container {padding: 0px 15px 0px 15px!important;mso-padding-alt: 0px 15px 0px 15px!important;}
			.mobile-column-left-padding { padding: 0px 0px 0px 0px!important; mso-alt-padding: 0px 0px 0px 0px!important; }
			.mobile-column-right-padding { padding: 0px 0px 0px 0px!important; mso-alt-padding: 0px 0px 0px 0px!important; }
			.mobile {display: none !important}
		  }
		  /**** My mobile styles ****/
		  @media only screen and (max-width: 599px) and (-webkit-min-device-pixel-ratio: 1) {
			*[class].wrapper { width:100% !important; }
			*[class].container { width:100% !important; }
			*[class].mobile { width:100% !important; display:block !important; }
			*[class].image{ width:100% !important; height:auto; }
			*[class].center{ margin:0 auto !important; text-align:center !important; }
			*[class="mobileOff"] { width: 0px !important; display: none !important; }
			*[class*="mobileOn"] { display: block !important; max-height:none !important; }
			p[class="mobile-padding"] {padding-left: 0px!important;padding-top: 10px;}
			.padding-container {padding: 0px 15px 0px 15px!important;mso-padding-alt: 0px 15px 0px 15px!important;}
			.hund {width: 100% !important;height: auto !important;}
			.td-padding {padding-left: 15px!important;padding-right: 15px!important;}
			.mobile-column-left-padding { padding: 18px 0px 18px 0px!important; mso-alt-padding: 18px 0px 18px 0px!important; }
			.mobile-column-right-padding { padding: 18px 0px 0px 0px!important; mso-alt-padding: 18px 0px 0px 0px!important; }
			.stack { width: 100% !important; }
			img {width: 100%!important;height: auto!important;}
			*[class="hide"] {display: none !important}
			*[class="Gmail"] {display: none !important}
			.Gmail {display: none !important}
			.bottom-padding-fix {padding: 0px 0px 18px 0px!important; mso-alt-padding: 0px 0px 18px 0px;}
		  }
		  .social, .social:active {
			opacity: 1!important;
			transform: scale(1);
			transition: all .2s!important;
		  }
		  .social:hover {
			opacity: 0.8!important;
			transform: scale(1.1);
			transition: all .2s!important;
		  }
		  .button.raised {
			transition: box-shadow 0.2s cubic-bezier(0.4, 0, 0.2, 1);
			transition: all .2s;box-shadow: 0 2px 5px 0 rgba(0, 0, 0, 0.26);
		  }
		  .button.raised:hover {
			box-shadow: 0 8px 17px 0 rgba(0, 0, 0, 0.2);transition: all .2s;
			-webkit-box-shadow: 0 8px 17px 0 rgba(0, 0, 0, 0.2);transition: all .2s;
			-moz-box-shadow: 0 8px 17px 0 rgba(0, 0, 0, 0.2);transition: all .2s;
		  }
		  .card-1 {
			box-shadow: 0 2px 5px 0 rgba(0,0,0,.16), 0 2px 10px 0 rgba(0,0,0,.12);
			-webkit-box-shadow: 0 2px 5px 0 rgba(0,0,0,.16), 0 2px 10px 0 rgba(0,0,0,.12);
			-moz-box-shadow: 0 2px 5px 0 rgba(0,0,0,.16), 0 2px 10px 0 rgba(0,0,0,.12);
			transition: box-shadow .45s;
		  }
		  .card-1:hover {
			box-shadow: 0 8px 17px 0 rgba(0,0,0,.2), 0 6px 20px 0 rgba(0,0,0,.19);
			-webkit-box-shadow: 0 8px 17px 0 rgba(0,0,0,.2), 0 6px 20px 0 rgba(0,0,0,.19);
			-moz-box-shadow: 0 8px 17px 0 rgba(0,0,0,.2), 0 6px 20px 0 rgba(0,0,0,.19);
			transition: box-shadow .45s;
		  }
		  .ripplelink{
			display:block
			color:#fff;
			text-decoration:none;
			position:relative;
			overflow:hidden;
			-webkit-transition: all 0.2s ease;
			-moz-transition: all 0.2s ease;
			-o-transition: all 0.2s ease;
			transition: all 0.2s ease;
			z-index:0;
		  }
		  .ripplelink:hover{
			z-index:1000;
		  }
		  .ink {
			display: block;
			position: absolute;
			background:rgba(255, 255, 255, 0.3);
			border-radius: 100%;
			-webkit-transform:scale(0);
			   -moz-transform:scale(0);
				 -o-transform:scale(0);
					transform:scale(0);
		  }
		  .animate {
			-webkit-animation:ripple 0.65s linear;
			 -moz-animation:ripple 0.65s linear;
			  -ms-animation:ripple 0.65s linear;
			   -o-animation:ripple 0.65s linear;
				  animation:ripple 0.65s linear;
		  }
		  @-webkit-keyframes ripple {
			  100% {opacity: 0; -webkit-transform: scale(2.5);}
		  }
		  @-moz-keyframes ripple {
			  100% {opacity: 0; -moz-transform: scale(2.5);}
		  }
		  @-o-keyframes ripple {
			  100% {opacity: 0; -o-transform: scale(2.5);}
		  }
		  @keyframes ripple {
			  100% {opacity: 0; transform: scale(2.5);}
		  }
	
		</style>
		<!--[if gte mso 9]>
		<xml>
		  <o:OfficeDocumentSettings>
			<o:AllowPNG/>
			<o:PixelsPerInch>96</o:PixelsPerInch>
		  </o:OfficeDocumentSettings>
		</xml>
		<![endif]-->
	  </head>
	  <body style="margin:0; padding:0; background-color: #eeeeee;" bgcolor="#eeeeee">
		<!--[if mso]>
		<style type="text/css">
		body, table, td {font-family: Arial, Helvetica, sans-serif !important;}
		</style>
		<![endif]-->
		<!-- START EMAIL -->
		<table width="100%" cellpadding="0" cellspacing="0" border="0" bgcolor="#eeeeee">
		  <div class="Gmail" style="height: 1px !important; margin-top: -1px !important; max-width: 600px !important; min-width: 600px !important; width: 600px !important;"></div>
		  <div style="display: none; max-height: 0px; overflow: hidden;">
			
		  </div>
		  <!-- Insert &zwnj;&nbsp; hack after hidden preview text -->
		  <div style="display: none; max-height: 0px; overflow: hidden;">
			&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
		  </div>
	
		  <!-- START LOGO -->
		  <tr>
			<td width="100%" valign="top" align="center" class="padding-container" style="padding: 18px 0px 18px 0px!important; mso-padding-alt: 18px 0px 18px 0px;">
			  <table width="600" cellpadding="0" cellspacing="0" border="0" align="center" class="wrapper">
				<tr>
				  <td align="center">
					<table cellpadding="0" cellspacing="0" border="0">
					  <tr>
						<td width="100%" valign="top" align="center">
						  <table width="600" cellpadding="0" cellspacing="0" border="0" align="center" class="wrapper" bgcolor="#eeeeee">
							<tr>
							  <td align="center">
								<table width="600" cellpadding="0" cellspacing="0" border="0" class="container" align="center">
								  <!-- START HEADER IMAGE -->
								  <tr>
									<td align="center" class="hund" width="600">
									  <img src="https://storage.googleapis.com/rails-2gaijin-storage/2gaijinheader.png" width="300" alt="Logo" border="0" style="max-width: 300px; display:block;">
									</td>
								  </tr>
								  <!-- END HEADER IMAGE -->
								</table>
							  </td>
							</tr>
						  </table>
						</td>
					  </tr>
					</table>
				  </td>
				</tr>
			  </table>
			</td>
		  </tr>
		  <!-- END LOGO -->
	
		  <!-- START CARD 1 -->
		  <tr>
			<td width="100%" valign="top" align="center" class="padding-container" style="padding-top: 0px!important; padding-bottom: 18px!important; mso-padding-alt: 0px 0px 18px 0px;">
			  <table width="600" cellpadding="0" cellspacing="0" border="0" align="center" class="wrapper">
				<tr>
				  <td>
					<table cellpadding="0" cellspacing="0" border="0">
					  <tr>
						<td style="border-radius: 3px; border-bottom: 2px solid #d4d4d4;" class="card-1" width="100%" valign="top" align="center">
						  <table style="border-radius: 3px;" width="600" cellpadding="0" cellspacing="0" border="0" align="center" class="wrapper" bgcolor="#ffffff">
							<tr>
							  <td align="center">
								<table width="600" cellpadding="0" cellspacing="0" border="0" class="container">
								  <!-- START HEADER IMAGE -->
								  <tr>
									<td align="center" class="hund ripplelink" width="400" style="padding-top: 20px;">
									  <img align="center" width="200" style="border-radius: 3px 3px 0px 0px; width: 100%; max-width: 200px!important" class="hund" src="https://img.icons8.com/color/452/talk-male--v2.png">
									</td>
								  </tr>
								  <!-- END HEADER IMAGE -->
								  <!-- START BODY COPY -->
								  <tr>
									<td class="td-padding" align="left" style="font-family: 'Roboto Mono', monospace; color: #212121!important; font-size: 24px; line-height: 30px; padding-top: 18px; padding-left: 18px!important; padding-right: 18px!important; padding-bottom: 0px!important; mso-line-height-rule: exactly; mso-padding-alt: 18px 18px 0px 13px;">
									  Your order for ` + itemName + ` has been accepted!
									</td>
								  </tr>
								  <!-- END BODY COPY -->
								  <!-- BUTTON -->
								  <tr>
									<td align="center" style="padding: 18px 18px 18px 18px; mso-alt-padding: 18px 18px 18px 18px!important;">
									  <a class="button raised" href="https://webbeta06012020.2gaijin.com" target="_blank" style="font-size: 24px; line-height: 14px; font-weight: 500; font-family: Helvetica, Arial, sans-serif; background-color: #ff8c00; color: #ffffff; text-decoration: none; border-radius: 3px; padding: 25px 100px; border: 1px solid #17bef7; display: inline-block;">Click Here to Check your Order</a>
									  <!-- <table width="100%" border="0" cellspacing="0" cellpadding="0">
										<tr>
										  <td>
											<table border="0" cellspacing="0" cellpadding="0">
											  <tr>
												<td align="center" style="border-radius: 3px;" bgcolor="#17bef7">
												  
												</td>
											  </tr>
											</table>
										  </td>
										</tr>
									  </table> -->
									</td>
								  </tr>
								  <!-- END BUTTON -->
								</table>
							  </td>
							</tr>
						  </table>
						</td>
					  </tr>
					</table>
				  </td>
				</tr>
			  </table>
			</td>
		  </tr>
		  <!-- END CARD 1 -->
		  <!-- FOOTER -->
		  <tr>
			<td width="100%" valign="top" align="center" class="padding-container">
			  <table width="600" cellpadding="0" cellspacing="0" border="0" align="center" class="wrapper">
				<tr>
				  <td width="100%" valign="top" align="center">
					<table width="600" cellpadding="0" cellspacing="0" border="0" align="center" class="wrapper" bgcolor="#eeeeee">
					  <tr>
						<td align="center">
						  <table width="600" cellpadding="0" cellspacing="0" border="0" class="container">
							<tr>
							  <!-- SOCIAL -->
							  <td align="center" width="300" style="padding-top: 0px!important; padding-bottom: 18px!important; mso-padding-alt: 0px 0px 18px 0px;">
								<table border="0" cellspacing="0" cellpadding="0">
								  <tr>
									<td align="right" valign="top" class="social">
									  <a href="https://www.facebook.com/2gaijin/"
									  target="_blank">
									  <img src="http://paulgoddarddesign.com/emails/images/material-design/fb-icon.png"
									  height="24" alt="Facebook" border="0" style="display:block; max-width: 24px">
									  </a>
									</td>
									<td width="20"></td>
									<td align="right" valign="top" class="social">
									  <a href="mailto:2gaijin@kitalabs.com"
									  target="_blank">
									  <img src="https://cdn4.iconfinder.com/data/icons/new-google-logo-2015/400/new-google-favicon-512.png"
									  height="24" alt="Google" border="0" style="display:block; max-width: 24px">
									  </a>
									</td>
								  </tr>
								</table>
							  </td>
							  <!-- END SOCIAL -->
							</tr>
							<tr>
							  <td class="td-padding" align="center" style="font-family: 'Roboto Mono', monospace; color: #212121!important; font-size: 16px; line-height: 24px; padding-top: 0px; padding-left: 0px!important; padding-right: 0px!important; padding-bottom: 0px!important; mso-line-height-rule: exactly; mso-padding-alt: 0px 0px 0px 0px;">
								&copy; 2019 2Gaijin.com
							  </td>
							</tr>
						  </table>
						</td>
					  </tr>
					</table>
				  </td>
				</tr>
			  </table>
			</td>
		  </tr>
		  <!-- FOOTER -->
	
		  <!-- SPACER -->
		  <!--[if gte mso 9]>
		  <table align="center" border="0" cellspacing="0" cellpadding="0" width="600">
			<tr>
			  <td align="center" valign="top" width="600" height="36">
				<![endif]-->
				<tr><td height="36px"></td></tr>
				<!--[if gte mso 9]>
			  </td>
			</tr>
		  </table>
		  <![endif]-->
		  <!-- END SPACER -->
	
		</table>
		<!-- END EMAIL -->
		<div style="display:none; white-space:nowrap; font:15px courier; line-height:0;">
		  &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp;
		  &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp;
		  &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp;
		</div>
	  </body>
	</html>
	`
	return template
}

func OrderRejectedEmail(itemName string) string {
	template := `<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01 Transitional//EN" "http://www.w3.org/TR/html4/loose.dtd">
	<html xmlns="http://www.w3.org/1999/xhtml" xmlns:v="urn:schemas-microsoft-com:vml" xmlns:o="urn:schemas-microsoft-com:office:office">
	  <head>
		<meta content="text/html; charset=UTF-8" http-equiv="Content-Type" />
		<!-- [ if !mso]> <!-->
		<meta content="IE=edge" http-equiv="X-UA-Compatible" />
		<!-- <![endif] -->
		<meta content="telephone=no" name="format-detection" />
		<meta content="width=device-width, initial-scale=1.0" name="viewport" />
		<link rel="apple-touch-icon" sizes="76x76" href="https://storage.cloud.google.com/rails-2gaijin-storage/2gaijinicon.png">
		<link rel="icon" type="image/png" sizes="96x96" href="https://storage.cloud.google.com/rails-2gaijin-storage/2gaijinicon.png">
		<title>2Gaijin.com: We are migrating to new website!</title>
		<link href="https://fonts.googleapis.com/css?family=Roboto+Mono" rel="stylesheet">
		<script src='http://cdnjs.cloudflare.com/ajax/libs/jquery/2.1.3/jquery.min.js'></script>
		<script src="http://paulgoddarddesign.com/js/ripple.js"></script>
		<style type="text/css">
		  .ExternalClass {width: 100%;}
		  .ExternalClass, .ExternalClass p, .ExternalClass span, .ExternalClass font, .ExternalClass td, .ExternalClass div, .ExternalClass b, .ExternalClass br, .ExternalClass img {line-height: 100% !important;}
		  /* iOS BLUE LINKS */
		  .appleBody a {color:#212121; text-decoration: none;}
		  .appleFooter a {color:#212121!important; text-decoration: none!important;}
		  /* END iOS BLUE LINKS */
		  img {color: #ffffff;text-align: center;font-family: Open Sans, Helvetica, Arial, sans-serif;display: block;}
		  body {margin: 0;padding: 0;-webkit-text-size-adjust: 100% !important;-ms-text-size-adjust: 100% !important;font-family: 'Open Sans', Helvetica, Arial, sans-serif!important;}
		  body,#body_style {background: #fffffe;}
		  table td {border-collapse: collapse;border-spacing: 0 !important;}
		  table tr {border-collapse: collapse;border-spacing: 0 !important;}
		  table tbody {border-collapse: collapse;border-spacing: 0 !important;}
		  table {border-collapse: collapse;border-spacing: 0 !important;}
		  span.yshortcuts,a span.yshortcuts {color: #000001;background-color: none;border: none;}
		  span.yshortcuts:hover,
		  span.yshortcuts:active,
		  span.yshortcuts:focus {color: #000001; background-color: none; border: none;}
		  img {-ms-interpolation-mode: : bicubic;}
		  a[x-apple-data-detectors] {color: inherit !important;text-decoration: none !important;font-size: inherit !important;font-family: inherit !important;font-weight: inherit !important;line-height: inherit !important;
		  }
		  /**** My desktop styles ****/
		  @media only screen and (min-width: 600px) {
			.noDesk {display: none !important;}
			.td-padding {padding-left: 15px!important;padding-right: 15px!important;}
			.padding-container {padding: 0px 15px 0px 15px!important;mso-padding-alt: 0px 15px 0px 15px!important;}
			.mobile-column-left-padding { padding: 0px 0px 0px 0px!important; mso-alt-padding: 0px 0px 0px 0px!important; }
			.mobile-column-right-padding { padding: 0px 0px 0px 0px!important; mso-alt-padding: 0px 0px 0px 0px!important; }
			.mobile {display: none !important}
		  }
		  /**** My mobile styles ****/
		  @media only screen and (max-width: 599px) and (-webkit-min-device-pixel-ratio: 1) {
			*[class].wrapper { width:100% !important; }
			*[class].container { width:100% !important; }
			*[class].mobile { width:100% !important; display:block !important; }
			*[class].image{ width:100% !important; height:auto; }
			*[class].center{ margin:0 auto !important; text-align:center !important; }
			*[class="mobileOff"] { width: 0px !important; display: none !important; }
			*[class*="mobileOn"] { display: block !important; max-height:none !important; }
			p[class="mobile-padding"] {padding-left: 0px!important;padding-top: 10px;}
			.padding-container {padding: 0px 15px 0px 15px!important;mso-padding-alt: 0px 15px 0px 15px!important;}
			.hund {width: 100% !important;height: auto !important;}
			.td-padding {padding-left: 15px!important;padding-right: 15px!important;}
			.mobile-column-left-padding { padding: 18px 0px 18px 0px!important; mso-alt-padding: 18px 0px 18px 0px!important; }
			.mobile-column-right-padding { padding: 18px 0px 0px 0px!important; mso-alt-padding: 18px 0px 0px 0px!important; }
			.stack { width: 100% !important; }
			img {width: 100%!important;height: auto!important;}
			*[class="hide"] {display: none !important}
			*[class="Gmail"] {display: none !important}
			.Gmail {display: none !important}
			.bottom-padding-fix {padding: 0px 0px 18px 0px!important; mso-alt-padding: 0px 0px 18px 0px;}
		  }
		  .social, .social:active {
			opacity: 1!important;
			transform: scale(1);
			transition: all .2s!important;
		  }
		  .social:hover {
			opacity: 0.8!important;
			transform: scale(1.1);
			transition: all .2s!important;
		  }
		  .button.raised {
			transition: box-shadow 0.2s cubic-bezier(0.4, 0, 0.2, 1);
			transition: all .2s;box-shadow: 0 2px 5px 0 rgba(0, 0, 0, 0.26);
		  }
		  .button.raised:hover {
			box-shadow: 0 8px 17px 0 rgba(0, 0, 0, 0.2);transition: all .2s;
			-webkit-box-shadow: 0 8px 17px 0 rgba(0, 0, 0, 0.2);transition: all .2s;
			-moz-box-shadow: 0 8px 17px 0 rgba(0, 0, 0, 0.2);transition: all .2s;
		  }
		  .card-1 {
			box-shadow: 0 2px 5px 0 rgba(0,0,0,.16), 0 2px 10px 0 rgba(0,0,0,.12);
			-webkit-box-shadow: 0 2px 5px 0 rgba(0,0,0,.16), 0 2px 10px 0 rgba(0,0,0,.12);
			-moz-box-shadow: 0 2px 5px 0 rgba(0,0,0,.16), 0 2px 10px 0 rgba(0,0,0,.12);
			transition: box-shadow .45s;
		  }
		  .card-1:hover {
			box-shadow: 0 8px 17px 0 rgba(0,0,0,.2), 0 6px 20px 0 rgba(0,0,0,.19);
			-webkit-box-shadow: 0 8px 17px 0 rgba(0,0,0,.2), 0 6px 20px 0 rgba(0,0,0,.19);
			-moz-box-shadow: 0 8px 17px 0 rgba(0,0,0,.2), 0 6px 20px 0 rgba(0,0,0,.19);
			transition: box-shadow .45s;
		  }
		  .ripplelink{
			display:block
			color:#fff;
			text-decoration:none;
			position:relative;
			overflow:hidden;
			-webkit-transition: all 0.2s ease;
			-moz-transition: all 0.2s ease;
			-o-transition: all 0.2s ease;
			transition: all 0.2s ease;
			z-index:0;
		  }
		  .ripplelink:hover{
			z-index:1000;
		  }
		  .ink {
			display: block;
			position: absolute;
			background:rgba(255, 255, 255, 0.3);
			border-radius: 100%;
			-webkit-transform:scale(0);
			   -moz-transform:scale(0);
				 -o-transform:scale(0);
					transform:scale(0);
		  }
		  .animate {
			-webkit-animation:ripple 0.65s linear;
			 -moz-animation:ripple 0.65s linear;
			  -ms-animation:ripple 0.65s linear;
			   -o-animation:ripple 0.65s linear;
				  animation:ripple 0.65s linear;
		  }
		  @-webkit-keyframes ripple {
			  100% {opacity: 0; -webkit-transform: scale(2.5);}
		  }
		  @-moz-keyframes ripple {
			  100% {opacity: 0; -moz-transform: scale(2.5);}
		  }
		  @-o-keyframes ripple {
			  100% {opacity: 0; -o-transform: scale(2.5);}
		  }
		  @keyframes ripple {
			  100% {opacity: 0; transform: scale(2.5);}
		  }
	
		</style>
		<!--[if gte mso 9]>
		<xml>
		  <o:OfficeDocumentSettings>
			<o:AllowPNG/>
			<o:PixelsPerInch>96</o:PixelsPerInch>
		  </o:OfficeDocumentSettings>
		</xml>
		<![endif]-->
	  </head>
	  <body style="margin:0; padding:0; background-color: #eeeeee;" bgcolor="#eeeeee">
		<!--[if mso]>
		<style type="text/css">
		body, table, td {font-family: Arial, Helvetica, sans-serif !important;}
		</style>
		<![endif]-->
		<!-- START EMAIL -->
		<table width="100%" cellpadding="0" cellspacing="0" border="0" bgcolor="#eeeeee">
		  <div class="Gmail" style="height: 1px !important; margin-top: -1px !important; max-width: 600px !important; min-width: 600px !important; width: 600px !important;"></div>
		  <div style="display: none; max-height: 0px; overflow: hidden;">
			
		  </div>
		  <!-- Insert &zwnj;&nbsp; hack after hidden preview text -->
		  <div style="display: none; max-height: 0px; overflow: hidden;">
			&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
		  </div>
	
		  <!-- START LOGO -->
		  <tr>
			<td width="100%" valign="top" align="center" class="padding-container" style="padding: 18px 0px 18px 0px!important; mso-padding-alt: 18px 0px 18px 0px;">
			  <table width="600" cellpadding="0" cellspacing="0" border="0" align="center" class="wrapper">
				<tr>
				  <td align="center">
					<table cellpadding="0" cellspacing="0" border="0">
					  <tr>
						<td width="100%" valign="top" align="center">
						  <table width="600" cellpadding="0" cellspacing="0" border="0" align="center" class="wrapper" bgcolor="#eeeeee">
							<tr>
							  <td align="center">
								<table width="600" cellpadding="0" cellspacing="0" border="0" class="container" align="center">
								  <!-- START HEADER IMAGE -->
								  <tr>
									<td align="center" class="hund" width="600">
									  <img src="https://storage.googleapis.com/rails-2gaijin-storage/2gaijinheader.png" width="300" alt="Logo" border="0" style="max-width: 300px; display:block;">
									</td>
								  </tr>
								  <!-- END HEADER IMAGE -->
								</table>
							  </td>
							</tr>
						  </table>
						</td>
					  </tr>
					</table>
				  </td>
				</tr>
			  </table>
			</td>
		  </tr>
		  <!-- END LOGO -->
	
		  <!-- START CARD 1 -->
		  <tr>
			<td width="100%" valign="top" align="center" class="padding-container" style="padding-top: 0px!important; padding-bottom: 18px!important; mso-padding-alt: 0px 0px 18px 0px;">
			  <table width="600" cellpadding="0" cellspacing="0" border="0" align="center" class="wrapper">
				<tr>
				  <td>
					<table cellpadding="0" cellspacing="0" border="0">
					  <tr>
						<td style="border-radius: 3px; border-bottom: 2px solid #d4d4d4;" class="card-1" width="100%" valign="top" align="center">
						  <table style="border-radius: 3px;" width="600" cellpadding="0" cellspacing="0" border="0" align="center" class="wrapper" bgcolor="#ffffff">
							<tr>
							  <td align="center">
								<table width="600" cellpadding="0" cellspacing="0" border="0" class="container">
								  <!-- START HEADER IMAGE -->
								  <tr>
									<td align="center" class="hund ripplelink" width="400" style="padding-top: 20px;">
									  <img align="center" width="200" style="border-radius: 3px 3px 0px 0px; width: 100%; max-width: 200px!important" class="hund" src="https://img.icons8.com/color/452/talk-male--v2.png">
									</td>
								  </tr>
								  <!-- END HEADER IMAGE -->
								  <!-- START BODY COPY -->
								  <tr>
									<td class="td-padding" align="left" style="font-family: 'Roboto Mono', monospace; color: #212121!important; font-size: 24px; line-height: 30px; padding-top: 18px; padding-left: 18px!important; padding-right: 18px!important; padding-bottom: 0px!important; mso-line-height-rule: exactly; mso-padding-alt: 18px 18px 0px 13px;">
									  Your order for ` + itemName + ` is rejected
									</td>
								  </tr>
								  <!-- END BODY COPY -->
								  <!-- BUTTON -->
								  <tr>
									<td align="center" style="padding: 18px 18px 18px 18px; mso-alt-padding: 18px 18px 18px 18px!important;">
									  <a class="button raised" href="https://webbeta06012020.2gaijin.com" target="_blank" style="font-size: 24px; line-height: 14px; font-weight: 500; font-family: Helvetica, Arial, sans-serif; background-color: #ff8c00; color: #ffffff; text-decoration: none; border-radius: 3px; padding: 25px 100px; border: 1px solid #17bef7; display: inline-block;">Click Here to Check your Order</a>
									  <!-- <table width="100%" border="0" cellspacing="0" cellpadding="0">
										<tr>
										  <td>
											<table border="0" cellspacing="0" cellpadding="0">
											  <tr>
												<td align="center" style="border-radius: 3px;" bgcolor="#17bef7">
												  
												</td>
											  </tr>
											</table>
										  </td>
										</tr>
									  </table> -->
									</td>
								  </tr>
								  <!-- END BUTTON -->
								</table>
							  </td>
							</tr>
						  </table>
						</td>
					  </tr>
					</table>
				  </td>
				</tr>
			  </table>
			</td>
		  </tr>
		  <!-- END CARD 1 -->
		  <!-- FOOTER -->
		  <tr>
			<td width="100%" valign="top" align="center" class="padding-container">
			  <table width="600" cellpadding="0" cellspacing="0" border="0" align="center" class="wrapper">
				<tr>
				  <td width="100%" valign="top" align="center">
					<table width="600" cellpadding="0" cellspacing="0" border="0" align="center" class="wrapper" bgcolor="#eeeeee">
					  <tr>
						<td align="center">
						  <table width="600" cellpadding="0" cellspacing="0" border="0" class="container">
							<tr>
							  <!-- SOCIAL -->
							  <td align="center" width="300" style="padding-top: 0px!important; padding-bottom: 18px!important; mso-padding-alt: 0px 0px 18px 0px;">
								<table border="0" cellspacing="0" cellpadding="0">
								  <tr>
									<td align="right" valign="top" class="social">
									  <a href="https://www.facebook.com/2gaijin/"
									  target="_blank">
									  <img src="http://paulgoddarddesign.com/emails/images/material-design/fb-icon.png"
									  height="24" alt="Facebook" border="0" style="display:block; max-width: 24px">
									  </a>
									</td>
									<td width="20"></td>
									<td align="right" valign="top" class="social">
									  <a href="mailto:2gaijin@kitalabs.com"
									  target="_blank">
									  <img src="https://cdn4.iconfinder.com/data/icons/new-google-logo-2015/400/new-google-favicon-512.png"
									  height="24" alt="Google" border="0" style="display:block; max-width: 24px">
									  </a>
									</td>
								  </tr>
								</table>
							  </td>
							  <!-- END SOCIAL -->
							</tr>
							<tr>
							  <td class="td-padding" align="center" style="font-family: 'Roboto Mono', monospace; color: #212121!important; font-size: 16px; line-height: 24px; padding-top: 0px; padding-left: 0px!important; padding-right: 0px!important; padding-bottom: 0px!important; mso-line-height-rule: exactly; mso-padding-alt: 0px 0px 0px 0px;">
								&copy; 2019 2Gaijin.com
							  </td>
							</tr>
						  </table>
						</td>
					  </tr>
					</table>
				  </td>
				</tr>
			  </table>
			</td>
		  </tr>
		  <!-- FOOTER -->
	
		  <!-- SPACER -->
		  <!--[if gte mso 9]>
		  <table align="center" border="0" cellspacing="0" cellpadding="0" width="600">
			<tr>
			  <td align="center" valign="top" width="600" height="36">
				<![endif]-->
				<tr><td height="36px"></td></tr>
				<!--[if gte mso 9]>
			  </td>
			</tr>
		  </table>
		  <![endif]-->
		  <!-- END SPACER -->
	
		</table>
		<!-- END EMAIL -->
		<div style="display:none; white-space:nowrap; font:15px courier; line-height:0;">
		  &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp;
		  &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp;
		  &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp;
		</div>
	  </body>
	</html>
	`
	return template
}
