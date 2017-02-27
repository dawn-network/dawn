/*
 * Slider Settings
 */
jQuery(document).ready(function(){
   jQuery('.widget_slider_area_rotate').bxSlider({
      mode: 'horizontal',
      speed: 1500,
      auto: true,
      pause: 5000,
      adaptiveHeight: true,
      nextText: '',
      prevText: '',
      nextSelector: '.slide-next',
      prevSelector: '.slide-prev',
      pager: false,
      tickerHover: true
   });
});