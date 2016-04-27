/* FILLED WITH DUMMY DATA */

var greens = [
  "#0198E1", "#5A6351","#3B5323","#458B00","#7CFC00","#8BA870","#629632","#7F8778","#BCED91","#6A8455","#397D02","#567E3A","#A6D785","#687E5A","#8AA37B","#476A34","#93DB70","#4CBB17","#49E20E","#55AE3A","#C5E3BF","#A9C9A4","#8CDD81","#5F755E","#4AC948","#698B69","#548B54","#9BCD9B","#215E21","#B4EEB4","#003300","#008000","#00CD00","#33FF33","#C1FFC1","#4BB74C","#B2D0B4","#5B9C64","#0AC92B","#31B94D","#96C8A2","#00611C","#91B49C","#337147","#40664D","#0E8C3A","#2C5D3F","#00FF66","#00C957","#34925E","#426352","#006633","#00EE76","#2E473B","#213D30","#4D7865","#00FA9A","#458B74","#32CD99","#7FFFD4","#138F6A","#3B8471","#284942","#668E86","#174038","#4A766E","#3E766D","#DBFEF8","#2A8E82","#40E0D0","#03A89E","#01C5BB","#068481","#636F57","#78AB46","#66CD00","#7FFF00","#9CBA7F","#659D32","#586949","#488214","#CFDBC5","#748269","#9DB68C","#61B329","#3F602B","#8FA880","#5DFC0A","#3B5E2B","#484D46","#308014","#7BBF6A","#395D33","#86C67C","#699864","#B7C8B6","#838B83","#2F4F2F","#8FBC8F","#71C671","#228B22","#90EE90","#004F00","#008B00","#00EE00","#66FF66","#CCFFCC","#6EFF70","#24D330","#3F9E4D","#BDFCC9","#1DA237","#3EA055","#00AF33","#79A888","#37BC61","#92CCA6","#78A489","#759B84","#2E8B57","#54FF9F","#3CB371","#607C6E","#008B45","#00FF7F","#5EDA9E","#3E7A5E","#597368","#238E68","#66CDAA","#76EEC6","#218868","#00C78C","#49E9BD","#0FDDAF","#527F76","#20BF9F","#4F8E83","#4CB7A5","#DAF4F0","#353F3E","#99CDC9","#457371","#48D1CC","#526F35","#4A7023","#76EE00","#7F9A65","#3D5229","#324F17","#608341","#46523C","#577A3A","#83F52C","#C0D9AF","#77896C","#3F6826","#646F5E","#435D36","#84BE6A","#3A6629","#9CA998","#4DBD33","#596C56","#7BCC70","#3D8B37","#63AB62","#C1CDC1","#426F42","#E0EEE0","#7CCD7C","#32CD32","#98FB98","#006400","#009900","#00FF00","#9AFF9A","#F0FFF0","#3D9140","#4D6B50","#517B58","#3D5B43","#00FF33","#487153","#688571","#B4D7BF","#70DB93","#329555","#3E6B4F","#2E6444","#4EEE94","#DBE6E0","#43CD80","#43D58C","#00CD66","#F5FFFA","#B6C5BE","#327556","#28AE7B","#4C7064","#32CC99","#00FFAA","#808A87","#D0FAEE","#1B6453","#006B54","#A4DCD1","#00FFCC","#3E766D","#2FAA96","#108070","#36DBCA","#45C3B8","#20B2AA","#90FEFB"
]

var data01 = {
    labels: ["Free Space", "Volume01", "Volume03"],
    datasets: [
      {
          data: [300, 50, 100],
          backgroundColor: greens,
          hoverBackgroundColor: greens
      }]
};

var data02 = {
    labels: ["Free Space", "Volume01", "Volume03"],
    datasets: [
      {
          data: [300, 50, 100],
          backgroundColor: greens,
          hoverBackgroundColor: greens
      }]
};

var data03 = {
    labels: ["Free Space", "VPDX01", "PDX03"],
    datasets: [
      {
          data: [30, 500, 1000],
          backgroundColor: greens,
          hoverBackgroundColor: greens
      }]
};

var data04 = {
    labels: ["Free Space", "TEST01", "TEST03", "TEST04", "TEST05"],
    datasets: [
      {
          data: [300, 502, 780, 234, 904, 343],
          backgroundColor: greens,
          hoverBackgroundColor: greens
      }]
};

var data05 = {
    labels: ["Free Space", "Volume01", "Volume03"],
    datasets: [
      {
          data: [300, 503, 900],
          backgroundColor: greens,
          hoverBackgroundColor: greens
      }]
};

var options = {}



$(document).ready(function() {

  //call API to get all Storage Providers then start making a ton of stuff!

  var ss1 = document.getElementById("ss1");
  var ss1DoughnutChart = new Chart(ss1, {
    type: 'doughnut',
    data: data01,
    options: options
  });

  var ss2 = document.getElementById("ss2");
  var ss2DoughnutChart = new Chart(ss2, {
    type: 'doughnut',
    data: data02,
    options: options
  });

  var ss3 = document.getElementById("ss3");
  var ss3DoughnutChart = new Chart(ss3, {
    type: 'doughnut',
    data: data03,
    options: options
  });

  var ss4 = document.getElementById("ss4");
  var ss4DoughnutChart = new Chart(ss4, {
    type: 'doughnut',
    data: data04,
    options: options
  });

  var ss5 = document.getElementById("ss5");
  var ss5DoughnutChart = new Chart(ss5, {
    type: 'doughnut',
    data: data05,
    options: options
  });

});


Chart.defaults.global.legend.display = false;