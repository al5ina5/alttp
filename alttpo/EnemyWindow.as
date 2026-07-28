EnemyWindow @enemyWindow;

const string enemyHeader = "IDX HST HP   X    Y    PTR   ACT";

class EnemyWindow {
  GUI::Color clrYellow;
  GUI::Color clrEnabled;
  GUI::Color clrDisabled;
  GUI::Color clrWhite;

  GUI::Window @window;
  array<SM_Enemy> enemies(16);
  array<GUI::Label@> labels(16);

  EnemyWindow() {
    clrYellow   = GUI::Color(240, 240,   0);
    clrEnabled  = GUI::Color(200, 200, 200);
    clrDisabled = GUI::Color( 80,  80,  80);
    clrWhite    = GUI::Color(255, 255, 255);

    // relative position to bsnes window (top-right area)
    @window = GUI::Window(sx(300), sy(0), true);
    window.title = "SM Enemies";
    window.font = GUI::Font("{mono}", 8);
    window.size = GUI::Size(sx(44*8), sy(18*16));
    window.backgroundColor = GUI::Color(  0,   0,   0);

    auto @vl = GUI::VerticalLayout();
    vl.setPadding(sy(2), sy(2));

    // header row
    {
      auto @hl = GUI::HorizontalLayout();
      auto @lbl = GUI::Label();
      lbl.foregroundColor = clrYellow;
      lbl.text = enemyHeader;
      hl.append(lbl, GUI::Size(-1, 0));
      vl.append(hl, GUI::Size(-1, 0));
    }

    // divider
    {
      auto @div = GUI::Label();
      div.foregroundColor = clrYellow;
      div.text = "--------------------------------------------";
      vl.append(div, GUI::Size(-1, 0));
    }

    // 16 enemy rows
    for (uint i = 0; i < 16; i++) {
      auto @hl = GUI::HorizontalLayout();
      @labels[i] = GUI::Label();
      labels[i].foregroundColor = clrDisabled;
      labels[i].text = "";
      hl.append(labels[i], GUI::Size(-1, 0));
      vl.append(hl, GUI::Size(-1, 0));
    }

    window.append(vl);
    vl.resize();
    window.visible = true;
  }

  void show() {
    window.visible = true;
  }

  void hide() {
    window.visible = false;
  }

  // must be run from post_frame() to reflect current RAM state
  void update() {
    for (uint i = 0; i < 16; i++) {
      enemies[i].fetch(uint8(i));
      auto @en = enemies[i];

      if (en.is_active) {
        labels[i].foregroundColor = clrEnabled;
        labels[i].text =
          fmtHex(en.enemy_index, 2) + "  " +
          fmtHex(en.host_index, 2) + " " +
          "Y   " +
          fmtHex(en.health, 4) + " " +
          fmtHex(en.Xpos, 4) + " " +
          fmtHex(en.Ypos, 4) + " " +
          fmtHex(en.pointer, 4) + "  " +
          "1";
      } else {
        labels[i].foregroundColor = clrDisabled;
        labels[i].text =
          fmtHex(en.enemy_index, 2) + "  " +
          "-- -- ---- ---- ---- ----  0";
      }
    }
  }
};
