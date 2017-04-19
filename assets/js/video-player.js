var progressMonitorInterval;
var videoPlayer

$(document).ready(function() {
  videoPlayer = videojs('cohesion-video');
  videoPlayer.on('play', function(){
    console.log('playing')
    if(progressMonitorInterval == null){
      progressMonitorInterval = setInterval(monitorProgress, 1000)
    }
  })

  videoPlayer.on('pause', function(){
    console.log('paused')
  })

  videoPlayer.ready(function() {
    //  videoPlayer.currentTime(35);
  });
});


function monitorProgress(){
  if(videoPlayer.ended()){
    clearInterval(progressMonitorInterval)
    return
  }

  if(!videoPlayer.paused()){
    console.log(videoPlayer.currentTime());
  }
}
