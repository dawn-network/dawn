/*
 * Settings of the ticker
 */
jQuery(document).ready(function(){
   jQuery('.newsticker').newsTicker({
      row_height: 20,
      max_rows: 1,
      speed: 1000,
      direction: 'down',
      duration: 4000,
      autostart: 1,
      pauseOnHover: 1
   });
});