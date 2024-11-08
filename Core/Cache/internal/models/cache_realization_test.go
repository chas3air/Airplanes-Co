package models

import (
	"reflect"
	"testing"
	"time"
)

func TestCarrotCache_Get(t *testing.T) {
	cc := NewCarrotCache()
	cc.Set("testKey", "testValue", 10)

	tests := []struct {
		name  string
		args  string
		want  interface{}
		want1 bool
	}{
		{
			name:  "ExistingKey",
			args:  "testKey",
			want:  "testValue",
			want1: true,
		},
		{
			name:  "NonExistingKey",
			args:  "nonExistentKey",
			want:  nil,
			want1: false,
		},
	}

	// Run tests for existing and non-existing key before expiration
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := cc.Get(tt.args)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CarrotCache.Get() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("CarrotCache.Get() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}

	// Sleep to test expiration
	time.Sleep(11 * time.Second)

	// Now test the expired key
	expiredKeyTest := struct {
		name  string
		args  string
		want  interface{}
		want1 bool
	}{
		name:  "ExpiredKey",
		args:  "testKey",
		want:  nil,
		want1: false,
	}

	t.Run(expiredKeyTest.name, func(t *testing.T) {
		got, got1 := cc.Get(expiredKeyTest.args)
		if !reflect.DeepEqual(got, expiredKeyTest.want) {
			t.Errorf("CarrotCache.Get() got = %v, want %v", got, expiredKeyTest.want)
		}
		if got1 != expiredKeyTest.want1 {
			t.Errorf("CarrotCache.Get() got1 = %v, want %v", got1, expiredKeyTest.want1)
		}
	})
}

func TestCarrotCache_Set(t *testing.T) {
	cc := NewCarrotCache()

	tests := []struct {
		name string
		cc   *CarrotCache
		args struct {
			key   string
			value interface{}
			ttl   int
		}
	}{
		{
			name: "SetValue",
			cc:   cc,
			args: struct {
				key   string
				value interface{}
				ttl   int
			}{
				key:   "key1",
				value: "value1",
				ttl:   5,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.cc.Set(tt.args.key, tt.args.value, tt.args.ttl)
			got, found := tt.cc.Get(tt.args.key)
			if !found || !reflect.DeepEqual(got, tt.args.value) {
				t.Errorf("CarrotCache.Set() failed to store value, got = %v, want %v", got, tt.args.value)
			}
		})
	}
}

func TestCarrotCache_Delete(t *testing.T) {
	cc := NewCarrotCache()
	cc.Set("deleteKey", "deleteValue", 10)

	tests := []struct {
		name string
		cc   *CarrotCache
		args string
		want interface{}
	}{
		{
			name: "DeleteExistingKey",
			cc:   cc,
			args: "deleteKey",
			want: CacheItem{
				Value:      "deleteValue",
				Expiration: 10,
				SetedTime:  cc.cache["deleteKey"].SetedTime,
			},
		},
		{
			name: "DeleteNonExistingKey",
			cc:   cc,
			args: "nonExistentKey",
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cc.Delete(tt.args)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CarrotCache.Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCarrotCache_cleanUp(t *testing.T) {
	cc := NewCarrotCache()
	cc.Set("tempKey", "tempValue", 1) // Set a key that will expire quickly

	time.Sleep(2 * time.Second) // Wait for the key to expire
	cc.cleanUp()                // Trigger cleanup

	if _, found := cc.Get("tempKey"); found {
		t.Errorf("Expected 'tempKey' to be cleaned up, but it still exists")
	}
}
