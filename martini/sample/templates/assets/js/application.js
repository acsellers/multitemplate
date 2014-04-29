$(document).ready(function(){
  $(".key-link").click(function(ev){
    key_val = ev.currentTarget.dataset["key"];
    $(".seperator").hide();
    $(".keycontent").hide();
    $(".key-link").removeClass("active");
    $('a[data-key="'+key_val+'"]').addClass("active");
    $("#key_"+key_val).show();
  });

  $(".all_keys").click(function(){
    $(".seperator").show();
    $(".keycontent").show();
    $(".key-link").removeClass("active");
  });

  $(".redis-item").each(function(i,v){
    try {
      value = JSON.parse($(v).html());
      $(v).html(JSON.stringify(value, null, "    "));
    }
    catch(e) {
    }
  });
})
