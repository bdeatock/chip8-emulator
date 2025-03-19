let wasmReady = false;
window.addEventListener("message", function (event) {
  if (event.data && event.data.type === "wasmReady") {
    wasmReady = true;
  }
});

function refocusEmulator() {
  const iframe = document.querySelector("iframe");
  if (!iframe) return;

  iframe.focus();

  iframe.contentWindow.postMessage(
    {
      type: "focus",
    },
    window.location.origin
  );
}

function handleBodyClick(event) {
  if (
    event.target === document.body ||
    event.target.classList.contains("sidebar")
  ) {
    // Only refocus if click is directly on body or empty sidebar
    refocusEmulator();
  }
}

function handleFileSelect(event) {
  if (!wasmReady) return;

  const iframe = document.querySelector("iframe");
  if (!iframe) return;

  const file = event.target.files[0];
  if (!file) return;

  const reader = new FileReader();
  reader.onload = function (e) {
    const arrayBuffer = e.target.result;
    iframe.contentWindow.postMessage(
      {
        type: "loadROM",
        data: arrayBuffer,
      },
      window.location.origin
    );
    document.getElementById("rom-picker").value = "empty";
  };
  reader.readAsArrayBuffer(file);
}

function handleRomSelect(event) {
  if (!wasmReady) return;

  const iframe = document.querySelector("iframe");
  if (!iframe) return;

  const rom = event.target.value;
  if (!rom) return;

  fetch(`/roms/${rom}`)
    .then((response) => {
      if (!response.ok) {
        throw new Error(
          `Failed to load ROM: ${response.status} ${response.statusText}`
        );
      }
      return response.arrayBuffer();
    })
    .then((arrayBuffer) => {
      iframe.contentWindow.postMessage(
        {
          type: "loadROM",
          data: arrayBuffer,
        },
        window.location.origin
      );
      displayRomInfo(rom);
      refocusEmulator();
    })
    .catch((error) => {
      console.error("Error loading ROM:", error);
    });
}

function displayRomInfo(name) {
  const container = document.getElementById("rom-info-container");
  const title = document.getElementById("rom-info-title");
  const blurb = document.getElementById("rom-info-blurb");
  const controls = document.getElementById("rom-info-controls");

  switch (name) {
    case "brix":
      container.classList.remove("hidden");
      title.textContent = "Brix";
      blurb.textContent =
        "Smash through bricks by rebounding the ball with your paddle in this classic arcade game.";
      controls.innerHTML = "Controls:<br />'Q'/'E' - Move left or right";
      break;
    case "invaders":
      container.classList.remove("hidden");
      title.textContent = "Invaders";
      blurb.textContent =
        "Shoot the alien invaders before they reach the bottom of the screen.";
      controls.innerHTML =
        "Controls:<br />'Q'/'E' - Move left or right<br />'W' - Shoot<br /><br />Press 'W' to start game on main menu.";
      break;
    case "merlin":
      container.classList.remove("hidden");
      title.textContent = "Merlin";
      blurb.textContent = "Test your memory by repeating the pattern.";
      controls.innerHTML = "Controls:<br />'QWAS' - represent the 4 squares.";
      break;
    case "tetris":
      container.classList.remove("hidden");
      title.textContent = "Tetris";
      blurb.textContent = "";
      controls.innerHTML =
        "Controls:<br />'Q' - rotate.<br />'W'/'E' - Move left or right<br />'A' - Drop quickly";
      break;
    default:
      container.classList.add("hidden");
  }
}

function handleResetEmulator(event) {
  if (!wasmReady) return;

  const iframe = document.querySelector("iframe");
  if (!iframe) return;

  iframe.contentWindow.resetEmulator();

  refocusEmulator();
}

function handleToggleLegacyShift(event) {
  if (!wasmReady) return;

  const iframe = document.querySelector("iframe");
  if (!iframe) return;

  if (iframe.contentWindow.toggleLegacyShift()) {
    event.target.classList.add("toggle-on");
  } else {
    event.target.classList.remove("toggle-on");
  }

  refocusEmulator();
}

function handleToggleLegacyJump(event) {
  if (!wasmReady) return;

  const iframe = document.querySelector("iframe");
  if (!iframe) return;

  if (iframe.contentWindow.toggleLegacyJump()) {
    event.target.classList.add("toggle-on");
  } else {
    event.target.classList.remove("toggle-on");
  }

  refocusEmulator();
}

function handleToggleLegacyStoreLoad(event) {
  if (!wasmReady) return;

  const iframe = document.querySelector("iframe");
  if (!iframe) return;

  if (iframe.contentWindow.toggleLegacyStoreLoad()) {
    event.target.classList.add("toggle-on");
  } else {
    event.target.classList.remove("toggle-on");
  }

  refocusEmulator();
}

function handleSwitchMode(event) {
  if (!wasmReady) return;

  const iframe = document.querySelector("iframe");
  if (!iframe) return;

  const button = event.currentTarget;
  const label = document.querySelector(`label[for="${button.id}"]`);
  const i = document.querySelector("#pause-step-btn i");
  const tooltip = document.querySelector("#pause-step-btn .tooltiptext");

  if (iframe.contentWindow.switchMode()) {
    // we are paused
    label.textContent = "Step Mode";
    button.classList.remove("play-mode");
    i.classList.remove("fa-pause");
    i.classList.add("fa-play");
    tooltip.textContent = "Enter Run Mode (continuous)";
  } else {
    label.textContent = "Run Mode";
    button.classList.add("play-mode");
    i.classList.add("fa-pause");
    i.classList.remove("fa-play");
    tooltip.textContent = "Enter Step Mode (space bar - step)";
  }

  refocusEmulator();
}

function handleSetCycleRate(event) {
  if (!wasmReady) return;

  const iframe = document.querySelector("iframe");
  if (!iframe) return;

  iframe.contentWindow.updateCycleRate(parseInt(event.target.value));
  refocusEmulator();
}
