// Discord Game SDK is optional; this bsnes-as build may not register discord::
// bindings. Keep a tiny stub so the rest of ALttPO compiles and runs offline.
namespace discord {
  bool enabled = false;
  int result = 0;

  void cartridge_loaded() {}
  void pre_frame() {}
}
