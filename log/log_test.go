package log

import (
	"context"
	"testing"
)

type testCommon interface {
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fail()
	FailNow()
	Failed() bool
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Helper()
	Log(args ...interface{})
	Logf(format string, args ...interface{})
	Name() string
	Skip(args ...interface{})
	SkipNow()
	Skipf(format string, args ...interface{})
	Skipped() bool
}

func assertEqual(t testCommon, a, b interface{}) {
	t.Helper()
	if a != b {
		t.Errorf("Not Equal. %++v %++v", a, b)
	}
}

func assertTrue(t testCommon, a bool) {
	t.Helper()
	if !a {
		t.Errorf("Not True %++v", a)
	}
}

func TestConfig_Equals(t *testing.T) {
	config := Config{
		Rotate: RotateConfig{
			FilePath:   "FilePath",
			Filename:   "Filename",
			MaxSize:    100,
			MaxBackups: 50,
			MaxAge:     200,
		},
		Level:   "debug",
		Console: true,
	}

	other := Config{
		Rotate: RotateConfig{
			FilePath:   "FilePath",
			Filename:   "Filename",
			MaxSize:    100,
			MaxBackups: 50,
			MaxAge:     200,
		},
		Level:   "debug",
		Console: true,
	}

	assertTrue(t, config.Equals(&other))
}

func TestMetaDataWithMap(t *testing.T) {
	opts := options{}
	metaDataWithMapOpt := MetaDataWithMap(map[string]string{"key1": "value1"})
	metaDataWithMapOpt(&opts)
	assertEqual(t, len(opts.m), 1)
	assertEqual(t, opts.m["key1"], "value1")
}

func TestMetadata(t *testing.T) {
	opts := options{
		m: map[string]string{},
	}
	metaDataOpt := Metadata("key2", "value2")
	metaDataOpt(&opts)
	assertEqual(t, len(opts.m), 1)
	assertEqual(t, opts.m["key2"], "value2")

}
func TestBuildAuditKeyValues(t *testing.T) {
	exceptResult := []string{KeyTraceID, "0000001",
		KeyParentID, "0000003",
		KeySpanID, "0000002",
		"user", "user01",
		"operator", "operator01",
		"request", "request data",
		"operatorLevel", "1",
		"logtype", "audit",
	}
	ctx := context.Background()
	ctx = context.WithValue(ctx, KeyTraceID, "0000001")
	ctx = context.WithValue(ctx, KeyParentID, "0000003")
	ctx = context.WithValue(ctx, KeySpanID, "0000002")
	buildResult := buildAuditKeyValues(ctx, "user01", "operator01", "1", []byte("request data"), nil)
	assertEqual(t, len(buildResult), len(exceptResult))
	for i, v := range buildResult {
		if exceptResult[i] != v.(string) {
			t.Errorf("Unexcept result:%v(index=%d), except value:%v", v, i, exceptResult[i])
		}
	}
}

func TestBuildTraceKeyValues(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, KeyTraceID, "0000001")
	ctx = context.WithValue(ctx, KeyParentID, "0000003")
	ctx = context.WithValue(ctx, KeySpanID, "0000002")

	buildResult := buildStandardKeyValues(ctx, nil)
	exceptResult := []string{KeyTraceID, "0000001", KeyParentID, "0000003", KeySpanID, "0000002"}
	assertEqual(t, len(buildResult), len(exceptResult))
	for i, v := range buildResult {
		if exceptResult[i] != v.(string) {
			t.Errorf("Unexcept result:%v(index=%d), except value:%v", v, i, exceptResult[i])
		}
	}
}

func TestCheckNil(t *testing.T) {
	assertEqual(t, checkNil(nil), "")
	assertEqual(t, checkNil("test"), "test")
}

func TestName(_ *testing.T) {
	Init()
	RegisterHookFunc(func(ctx context.Context) (context.Context, []interface{}) {
		moduleId := ctx.Value("moduleId")
		if moduleId != nil {
			return ctx, []interface{}{"moduleId", moduleId, KeyTraceID, "modify: xxxx-traceId-xxxx", KeyParentID, "modify: parentId",
				KeySpanID, "modify: spanId", KeyTopicID, "modify: topicId"}
		}
		return ctx, nil
	})
	ctx := context.WithValue(context.Background(), KeyTraceID, "traceId-xxx-xxx")
	ctx = context.WithValue(ctx, KeySpanID, "spanId-xxx-xxx")
	ctx = context.WithValue(ctx, KeyServiceID, "serviceId-xxx")
	ctx = context.WithValue(ctx, KeyParentID, "parentId-xxx")
	ctx = context.WithValue(ctx, "moduleId", "nihao")

	Errorf(ctx, "contains moduleId")
	Errorf(context.Background(), "don't contains moduleId")
	Errorw(ctx, "contains moduleId and param contains", "moduleId", "hhh")
}

func BenchmarkTestErrorf(b *testing.B) {
	config.Console = false
	config.Rotate.FilePath = "./"
	Init()
	RegisterHookFunc(func(ctx context.Context) (context.Context, []interface{}) {
		moduleId := ctx.Value("moduleId")
		if moduleId != nil {
			return ctx, []interface{}{"moduleId", moduleId, KeyTraceID, "modify: xxxx-traceId-xxxx", KeyParentID, "modify: parentId",
				KeySpanID, "modify: spanId", KeyTopicID, "modify: topicId"}
		}
		return ctx, nil
	})
	ctx := context.WithValue(context.Background(), KeyTraceID, "traceId-xxx-xxx")
	ctx = context.WithValue(ctx, KeySpanID, "spanId-xxx-xxx")
	ctx = context.WithValue(ctx, KeyServiceID, "serviceId-xxx")
	ctx = context.WithValue(ctx, KeyParentID, "parentId-xxx")
	ctx = context.WithValue(ctx, "moduleId", "nihao")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Errorf(ctx, "contains moduleId")
	}
}


func BenchmarkTestErrorfWithoutHook(b *testing.B) {
	config.Console = false
	config.Rotate.FilePath = "./"
	Init()
	ctx := context.WithValue(context.Background(), KeyTraceID, "traceId-xxx-xxx")
	ctx = context.WithValue(ctx, KeySpanID, "spanId-xxx-xxx")
	ctx = context.WithValue(ctx, KeyServiceID, "serviceId-xxx")
	ctx = context.WithValue(ctx, KeyParentID, "parentId-xxx")
	ctx = context.WithValue(ctx, "moduleId", "nihao")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Errorf(ctx, "contains moduleId")
	}
}