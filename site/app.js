const ROMS = {
  brix: {
    name: "Brix",
    blurb:
      "Smash through bricks by rebounding the ball with your paddle in this classic arcade game.",
    controls: "Controls:<br />'Q'/'E' - Move left or right",
  },
  invaders: {
    name: "Invaders",
    blurb:
      "Shoot the alien invaders before they reach the bottom of the screen.",
    controls:
      "Controls:<br />'Q'/'E' - Move left or right<br />'W' - Shoot<br /><br />Press 'W' to start game on main menu.",
  },
  merlin: {
    name: "Merlin",
    blurb: "Test your memory by repeating the pattern.",
    controls: "Controls:<br />'QWAS' - represent the 4 squares.",
  },
  tetris: {
    name: "Tetris",
    blurb: "",
    controls:
      "Controls:<br />'Q' - rotate.<br />'W'/'E' - Move left or right<br />'A' - Drop quickly",
  },
};

const elements = {
  iframe: null,
  romPicker: null,
  romInfo: {
    container: null,
    title: null,
    blurb: null,
    controls: null,
  },
};

function cacheElements() {
  try {
    elements.iframe = document.querySelector("iframe");
    elements.romPicker = document.getElementById("rom-picker");
    elements.romInfo.container = document.getElementById("rom-info-container");
    elements.romInfo.title = document.getElementById("rom-info-title");
    elements.romInfo.blurb = document.getElementById("rom-info-blurb");
    elements.romInfo.controls = document.getElementById("rom-info-controls");

    // Verify all required elements exist
    for (const [key, element] of Object.entries(elements)) {
      if (!element) {
        console.error(`Required element not found: ${key}`);
        return;
      }
    }
  } catch (error) {
    console.error("Failed to cache elements:", error);
  }
}

let wasmReady = false;
window.addEventListener("message", function (event) {
  if (event.data && event.data.type === "wasmReady") {
    wasmReady = true;
    cacheElements();
  }
});

function refocusEmulator() {
  if (!elements.iframe) return;

  elements.iframe.focus();

  elements.iframe.contentWindow.postMessage(
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
  if (!wasmReady || !elements.iframe || !elements.romPicker) return;

  const file = event.target.files[0];
  if (!file) return;

  const reader = new FileReader();
  reader.onload = function (e) {
    const arrayBuffer = e.target.result;
    elements.iframe.contentWindow.postMessage(
      {
        type: "loadROM",
        data: arrayBuffer,
      },
      window.location.origin
    );
    elements.romPicker.value = "empty";
    refocusEmulator();
  };
  reader.readAsArrayBuffer(file);
}

function handleRomSelect(event) {
  if (!wasmReady || !elements.iframe || !elements.romPicker) return;

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
      elements.iframe.contentWindow.postMessage(
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
  if (
    !elements.romInfo.container ||
    !elements.romInfo.title ||
    !elements.romInfo.blurb ||
    !elements.romInfo.controls
  )
    return;

  const rom = ROMS[name];
  if (!rom) {
    elements.romInfo.container.classList.add("hidden");
    return;
  }

  elements.romInfo.container.classList.remove("hidden");
  elements.romInfo.title.textContent = rom.name;
  elements.romInfo.blurb.textContent = rom.blurb;
  elements.romInfo.controls.innerHTML = rom.controls;
}

function handleResetEmulator(event) {
  if (!wasmReady || !elements.iframe) return;

  elements.iframe.contentWindow.resetEmulator();

  refocusEmulator();
}

function handleToggleLegacyShift(event) {
  if (!wasmReady || !elements.iframe) return;

  if (elements.iframe.contentWindow.toggleLegacyShift()) {
    event.target.classList.add("toggle-on");
  } else {
    event.target.classList.remove("toggle-on");
  }

  refocusEmulator();
}

function handleToggleLegacyJump(event) {
  if (!wasmReady || !elements.iframe) return;

  if (elements.iframe.contentWindow.toggleLegacyJump()) {
    event.target.classList.add("toggle-on");
  } else {
    event.target.classList.remove("toggle-on");
  }

  refocusEmulator();
}

function handleToggleLegacyStoreLoad(event) {
  if (!wasmReady || !elements.iframe) return;

  if (elements.iframe.contentWindow.toggleLegacyStoreLoad()) {
    event.target.classList.add("toggle-on");
  } else {
    event.target.classList.remove("toggle-on");
  }

  refocusEmulator();
}

function handleSwitchMode(event) {
  if (!wasmReady || !elements.iframe) return;

  const button = event.currentTarget;
  const label = document.querySelector(`label[for="${button.id}"]`);
  const i = document.querySelector("#pause-step-btn i");
  const tooltip = document.querySelector("#pause-step-btn .tooltiptext");

  if (elements.iframe.contentWindow.switchMode()) {
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
  if (!wasmReady || !elements.iframe) return;

  elements.iframe.contentWindow.updateCycleRate(parseInt(event.target.value));
  refocusEmulator();
}
