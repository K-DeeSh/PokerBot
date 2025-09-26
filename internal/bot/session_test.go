package bot

import "testing"

func TestSessionApplyValue(t *testing.T) {
	sess := NewSession()
	sess.Await = StepHand
	if err := sess.ApplyValue("Ah Kh"); err != nil {
		t.Fatalf("unexpected hand error: %v", err)
	}
	if len(sess.Request.Hand) != 2 {
		t.Fatalf("hand not stored")
	}

	sess.Await = StepPlayers
	if err := sess.ApplyValue("5"); err != nil {
		t.Fatalf("unexpected players error: %v", err)
	}
	if sess.Request.Players != 5 {
		t.Fatalf("players not stored")
	}

	sess.Await = StepTrials
	if err := sess.ApplyValue("3000"); err != nil {
		t.Fatalf("unexpected trials error: %v", err)
	}

	sess.Await = StepBoard
	if err := sess.ApplyValue("Qh Jh Td"); err != nil {
		t.Fatalf("unexpected board error: %v", err)
	}

	if !sess.HasRequiredFields() {
		t.Fatal("expected session to have required fields")
	}
}

func TestSessionApplyValueErrors(t *testing.T) {
	sess := NewSession()
	sess.Await = StepHand
	if err := sess.ApplyValue("Ah"); err == nil {
		t.Fatal("expected error for missing card")
	}

	sess.Await = StepPlayers
	if err := sess.ApplyValue("1"); err == nil {
		t.Fatal("expected error for players")
	}
}
