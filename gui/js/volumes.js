$( document ).ready(function() {
	$.getJSON("http://127.0.0.1:8080/admin/volumes", function( data ) {
	  $.each( data, function( key, val ) {
  		var vol = '<div class="col-xs-6 col-md-3">'+
	          '<h4><center>' + val.name + '</center></h4>'+
	          '<div class="col-md-4 center-block">'+
	            '<img src="../images/volume.png" width="150%" class="img-responsive center-block">'+
	          '</div>'+
	          '<div class="col-md-8">'+
	          	'<table width="100%">'+
		          	'<tr><td>ID:</td><td> ' + val.id + '</td></tr>' +
		          	'<tr><td>Size:</td><td> ' + val.size + '</td></tr>' +
		          	'<tr><td>Type:</td><td> ' + val.type + '</td></tr>' +
		          	'<tr><td>Service:</td><td> ' + val.serviceName + '</td></tr>' +
		          	'<tr><td>Provider:</td><td> ' + val.storageProviderName + '</td></tr>' +
		          	'<tr><td>AZ:</td><td> ' + val.availabilityZone + '</td></tr>' +
		          	'<tr><td>Scheduler:</td><td> ' + val.scheduler + '</td></tr>' +
				'</table>' +
	          '</div>'+
	        '</div>';
  		$( "#volumeList" ).append(vol);
	  });
	});
});