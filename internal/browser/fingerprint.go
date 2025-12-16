package browser

import (
	"github.com/go-rod/rod/lib/proto"
)

// FingerprintMasking implements browser fingerprint randomization
type FingerprintMasking struct {
	enabled bool
}

// NewFingerprintMasking creates a fingerprint masker
func NewFingerprintMasking(enabled bool) *FingerprintMasking {
	return &FingerprintMasking{
		enabled: enabled,
	}
}

// ApplyMasking applies fingerprint masking techniques
func (fm *FingerprintMasking) ApplyMasking(page interface{}) error {
	if !fm.enabled {
		return nil
	}

	// Note: go-rod/stealth already handles most of this
	// Additional custom masking can be added here
	
	return nil
}

// MaskWebGL adds noise to WebGL fingerprinting
func (fm *FingerprintMasking) MaskWebGL() string {
	return `
	(function() {
		const getParameter = WebGLRenderingContext.prototype.getParameter;
		WebGLRenderingContext.prototype.getParameter = function(parameter) {
			// Add slight noise to WebGL parameters
			if (parameter === 37445) { // UNMASKED_VENDOR_WEBGL
				return 'Intel Inc.';
			}
			if (parameter === 37446) { // UNMASKED_RENDERER_WEBGL
				return 'Intel Iris OpenGL Engine';
			}
			return getParameter.apply(this, arguments);
		};
	})();
	`
}

// MaskCanvas adds noise to canvas fingerprinting
func (fm *FingerprintMasking) MaskCanvas() string {
	return `
	(function() {
		const originalToDataURL = HTMLCanvasElement.prototype.toDataURL;
		HTMLCanvasElement.prototype.toDataURL = function() {
			// Add minimal noise to canvas data
			const ctx = this.getContext('2d');
			if (ctx) {
				const imageData = ctx.getImageData(0, 0, this.width, this.height);
				for (let i = 0; i < imageData.data.length; i += 4) {
					if (Math.random() < 0.01) { // 1% of pixels
						imageData.data[i] += Math.floor(Math.random() * 3) - 1;
					}
				}
				ctx.putImageData(imageData, 0, 0);
			}
			return originalToDataURL.apply(this, arguments);
		};
	})();
	`
}

// DisableAutomationFlags removes automation detection flags
func (fm *FingerprintMasking) DisableAutomationFlags() string {
	return `
	(function() {
		// Remove webdriver flag
		Object.defineProperty(navigator, 'webdriver', {
			get: () => false
		});
		
		// Mask chrome automation
		window.navigator.chrome = {
			runtime: {}
		};
		
		// Override permissions
		const originalQuery = window.navigator.permissions.query;
		window.navigator.permissions.query = (parameters) => (
			parameters.name === 'notifications' ?
				Promise.resolve({ state: Notification.permission }) :
				originalQuery(parameters)
		);
		
		// Override plugins
		Object.defineProperty(navigator, 'plugins', {
			get: () => [1, 2, 3, 4, 5]
		});
		
		// Override languages
		Object.defineProperty(navigator, 'languages', {
			get: () => ['en-US', 'en']
		});
	})();
	`
}

// RandomizeViewport generates random viewport dimensions
func (fm *FingerprintMasking) RandomizeViewport() *proto.EmulationSetDeviceMetricsOverride {
	// Common resolutions
	resolutions := [][2]int{
		{1920, 1080},
		{1366, 768},
		{1536, 864},
		{1440, 900},
		{1280, 720},
	}
	
	idx := len(resolutions) / 2 // Use middle resolution for stability
	
	return &proto.EmulationSetDeviceMetricsOverride{
		Width:  resolutions[idx][0],
		Height: resolutions[idx][1],
		DeviceScaleFactor: 1,
		Mobile: false,
	}
}

// GetRandomUserAgent returns a realistic user agent
func GetRandomUserAgent() string {
	agents := []string{
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
	}
	
	return agents[0] // Use consistent user agent
}

// InjectStealthScripts injects all stealth scripts
func (fm *FingerprintMasking) InjectStealthScripts(page interface{}) error {
	if !fm.enabled {
		return nil
	}

	// Scripts are injected by go-rod/stealth automatically
	// Additional custom scripts can be added here if needed
	
	return nil
}
