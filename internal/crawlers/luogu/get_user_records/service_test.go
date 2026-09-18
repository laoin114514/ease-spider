package getuserrecords

import (
	"fmt"
	"spider/internal/models"
	"testing"
)

// oldDataSet 造 n 条"库里已有"的记录号
func oldDataSet(n int) map[string]bool {
	m := make(map[string]bool, n)
	for i := 0; i < n; i++ {
		m[fmt.Sprintf("sub-%d", i)] = true
	}
	return m
}

// 判断下一步动作：多出来的行不能再触发全量爬取，还在测评的记录要算作已同步
func TestPlanAfterIncremental(t *testing.T) {
	cases := []struct {
		name     string
		dbRows   int
		inserted int
		pending  int
		total    int
		want     crawlPlan
	}{
		{name: "库内与洛谷一致", dbRows: 20, total: 20, want: planConsistent},
		{name: "增量补齐了缺失的记录", dbRows: 18, inserted: 2, total: 20, want: planConsistent},
		{name: "刚提交还在测评也算同步", dbRows: 18, inserted: 1, pending: 1, total: 20, want: planConsistent},
		{name: "库内比洛谷少需要全量核对", dbRows: 10, inserted: 1, total: 20, want: planMissing},
		{name: "库内比洛谷多只需告警不再全量", dbRows: 20, total: 18, want: planExtra},
		{name: "库内多且本轮有新增仍然不补", dbRows: 96, inserted: 0, total: 95, want: planExtra},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			user := models.LuoguUserDeliver{
				OldDataSet: oldDataSet(tc.dbRows),
				Count:      tc.inserted,
				Pending:    tc.pending,
			}
			if got := planAfterIncremental(user, tc.total); got != tc.want {
				t.Fatalf("planAfterIncremental = %v, want %v (sync=%d total=%d)", got, tc.want, syncedCount(user), tc.total)
			}
		})
	}
}

func TestSyncedCount(t *testing.T) {
	user := models.LuoguUserDeliver{OldDataSet: oldDataSet(3), Count: 2, Pending: 1}
	if got := syncedCount(user); got != 6 {
		t.Fatalf("syncedCount = %d, want 6", got)
	}
}

// 同一个洛谷 uid 绑在多行 user 上时只保留第一行，并给出告警
func TestDedupeByUID(t *testing.T) {
	users := []models.LuoguUserDeliver{
		{Uid: "1822311", RealName: "张天杰"},
		{Uid: "1822311", RealName: "张天杰1"},
		{Uid: "123456", RealName: "甲"},
		{Uid: "1822311", RealName: "张天杰2"},
	}

	warns := make([]string, 0, 2)
	unique := dedupeByUID(users, func(format string, args ...interface{}) {
		warns = append(warns, fmt.Sprintf(format, args...))
	})

	if len(unique) != 2 {
		t.Fatalf("len(unique) = %d, want 2 (%+v)", len(unique), unique)
	}
	if unique[0].RealName != "张天杰" || unique[0].Uid != "1822311" {
		t.Fatalf("保留的行不对: %+v", unique[0])
	}
	if unique[1].RealName != "甲" {
		t.Fatalf("第二行不对: %+v", unique[1])
	}
	if len(warns) != 2 {
		t.Fatalf("告警条数 = %d, want 2 (%v)", len(warns), warns)
	}
	for _, w := range warns {
		if !contains(w, "1822311") || !contains(w, "重复绑定") {
			t.Fatalf("告警内容不对: %s", w)
		}
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (func() bool {
		for i := 0; i+len(sub) <= len(s); i++ {
			if s[i:i+len(sub)] == sub {
				return true
			}
		}
		return false
	})()
}

// 难度越界（洛谷新增难度等级）时回退为 unknown，不能越界 panic
func TestLuoguDifficulty(t *testing.T) {
	if got := luoguDifficulty(0); got == "" || got == "unknown" {
		t.Fatalf("difficulty 0 -> %q", got)
	}
	if got := luoguDifficulty(-1); got != "unknown" {
		t.Fatalf("difficulty -1 -> %q, want unknown", got)
	}
	if got := luoguDifficulty(99); got != "unknown" {
		t.Fatalf("difficulty 99 -> %q, want unknown", got)
	}
}
