/*
 * Gallery Post Format Slider Settings
 */
jQuery(document).ready(function(){
   jQuery('.gallery-images').bxSlider({
      mode: 'fade',
      speed: 1500,
      auto: true,
      pause: 3000,
      adaptiveHeight: true,
      nextText: '',
      prevText: '',
      nextSelector: '.slide-next',
      prevSelector: '.slide-prev',
      pager: false
   });
});