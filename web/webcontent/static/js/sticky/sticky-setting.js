/*
 * Settings of the sticky menu
 */

jQuery(document).ready(function(){
   var wpAdminBar = jQuery('#wpadminbar');
   if (wpAdminBar.length) {
      jQuery('#site-navigation').sticky({topSpacing:wpAdminBar.height()});
   } else {
      jQuery('#site-navigation').sticky({topSpacing:0});
   }
});