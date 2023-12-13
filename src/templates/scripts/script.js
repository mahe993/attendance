const timeBox = document.querySelector("#time-box");

(function currTime() {
  if (timeBox === null) return;
  timeBox.innerText = formattedDate();
  setInterval(() => {
    timeBox.innerText = formattedDate();
  }, 100);
})();

function formattedDate() {
  const options = {
    weekday: "short",
    day: "numeric",
    month: "short",
    hour: "numeric",
    minute: "numeric",
    second: "numeric",
    hour12: true,
  };

  const date = new Date().toLocaleString("en-US", options);
  return date;
}
