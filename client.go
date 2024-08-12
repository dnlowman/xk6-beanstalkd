package beanstalkd

import (
	"fmt"
	"time"

	"github.com/beanstalkd/go-beanstalk"
)

type Client struct {
	conn         *beanstalk.Conn
	usedTube     *beanstalk.Tube
	watchedTubes []string
	tubeSet      *beanstalk.TubeSet
}

func NewClient(addr string) (*Client, error) {
	conn, err := beanstalk.Dial("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Beanstalkd: %w", err)
	}
	c := &Client{
		conn:         conn,
		usedTube:     &beanstalk.Tube{Conn: conn, Name: "default"},
		watchedTubes: []string{"default"},
	}
	c.tubeSet = beanstalk.NewTubeSet(conn, "default")
	return c, nil
}

func (c *Client) Use(tube string) error {
	c.usedTube = &beanstalk.Tube{Conn: c.conn, Name: tube}
	return nil
}

func (c *Client) Watch(tube string) error {
	for _, t := range c.watchedTubes {
		if t == tube {
			return nil // Already watching this tube
		}
	}
	c.watchedTubes = append(c.watchedTubes, tube)
	c.tubeSet = beanstalk.NewTubeSet(c.conn, c.watchedTubes...)
	return nil
}

func (c *Client) Ignore(tube string) error {
	if len(c.watchedTubes) <= 1 {
		return fmt.Errorf("cannot ignore the only watched tube")
	}
	for i, t := range c.watchedTubes {
		if t == tube {
			c.watchedTubes = append(c.watchedTubes[:i], c.watchedTubes[i+1:]...)
			c.tubeSet = beanstalk.NewTubeSet(c.conn, c.watchedTubes...)
			return nil
		}
	}
	return fmt.Errorf("tube %s is not being watched", tube)
}

func (c *Client) Put(data string, priority uint32, delay, ttr time.Duration) (uint64, error) {
	return c.usedTube.Put([]byte(data), priority, delay, ttr)
}

func (c *Client) Reserve(timeout time.Duration) (uint64, string, error) {
	id, body, err := c.tubeSet.Reserve(timeout)
	if err != nil {
		return 0, "", fmt.Errorf("failed to reserve job: %w", err)
	}
	return id, string(body), nil
}

func (c *Client) Delete(id uint64) error {
	err := c.conn.Delete(id)
	if err != nil {
		return fmt.Errorf("failed to delete job %d: %w", id, err)
	}
	return nil
}

func (c *Client) Release(id uint64, priority uint32, delay time.Duration) error {
	err := c.conn.Release(id, priority, delay)
	if err != nil {
		return fmt.Errorf("failed to release job %d: %w", id, err)
	}
	return nil
}

func (c *Client) Bury(id uint64, priority uint32) error {
	err := c.conn.Bury(id, priority)
	if err != nil {
		return fmt.Errorf("failed to bury job %d: %w", id, err)
	}
	return nil
}

func (c *Client) Kick(bound int) (int, error) {
	kicked, err := c.conn.Kick(bound)
	if err != nil {
		return 0, fmt.Errorf("failed to kick jobs: %w", err)
	}
	return kicked, nil
}

func (c *Client) Peek(id uint64) ([]byte, error) {
	body, err := c.conn.Peek(id)
	if err != nil {
		return nil, fmt.Errorf("failed to peek job %d: %w", id, err)
	}
	return body, nil
}

func (c *Client) Stats() (map[string]string, error) {
	stats, err := c.conn.Stats()
	if err != nil {
		return nil, fmt.Errorf("failed to get stats: %w", err)
	}
	return stats, nil
}

func (c *Client) StatsTube(tube string) (map[string]string, error) {
	t := &beanstalk.Tube{Conn: c.conn, Name: tube}
	stats, err := t.Stats()
	if err != nil {
		return nil, fmt.Errorf("failed to get stats for tube %s: %w", tube, err)
	}
	return stats, nil
}

func (c *Client) ListTubes() ([]string, error) {
	tubes, err := c.conn.ListTubes()
	if err != nil {
		return nil, fmt.Errorf("failed to list tubes: %w", err)
	}
	return tubes, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
