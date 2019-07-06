package jvm

import (
	"log"
	"os"

	"github.com/linuxuser586/cassandra/pkg/normalize"
	"github.com/linuxuser586/common/pkg/net"
)

var podIP = net.PodIP

// Args are the Java commons args for Cassandra and tools
func Args() []string {
	ip, err := podIP()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	var s []string
	s = append(s, "-ea")
	s = append(s, "-XX:+UseThreadPriorities")
	s = append(s, "-XX:ThreadPriorityPolicy=42")
	s = append(s, "-XX:+HeapDumpOnOutOfMemoryError")
	s = append(s, "-Xss256k")
	s = append(s, "-XX:StringTableSize=1000003")
	s = append(s, "-XX:+AlwaysPreTouch")
	s = append(s, "-XX:-UseBiasedLocking")
	s = append(s, "-XX:+UseTLAB")
	s = append(s, "-XX:+ResizeTLAB")
	s = append(s, "-XX:+UseNUMA")
	s = append(s, "-XX:+PerfDisableSharedMem")
	s = append(s, "-XX:+UseG1GC")
	s = append(s, "-XX:G1RSetUpdatingPauseTimePercent=5")
	s = append(s, "-XX:MaxGCPauseMillis=500")
	gc := os.Getenv("LOG_GC")
	if normalize.True(gc) {
		s = append(s, "-Xloggc:/proc/self/fd/1")
		s = append(s, "-XX:+PrintGCDetails")
		s = append(s, "-XX:+PrintGCDateStamps")
		s = append(s, "-XX:+PrintHeapAtGC")
		s = append(s, "-XX:+PrintTenuringDistribution")
		s = append(s, "-XX:+PrintGCApplicationStoppedTime")
		s = append(s, "-XX:+PrintPromotionFailure")
		s = append(s, "-XX:PrintFLSStatistics=1")
		s = append(s, "-XX:-UseGCLogFileRotation")
	}
	s = append(s, "-XX:OnOutOfMemoryError=kill -9 %p")
	s = append(s, "-javaagent:/usr/share/cassandra/lib/jamm-0.3.0.jar")
	s = append(s, "-Dcassandra.jmx.remote.port=7199")
	s = append(s, "-Djava.rmi.server.hostname="+ip)
	s = append(s, "-Dcom.sun.management.jmxremote.authenticate=false")
	s = append(s, "-Dcom.sun.management.jmxremote.rmi.port=7199")
	s = append(s, "-Dcom.sun.management.jmxremote.password.file=/conf/jmxremote.password")
	s = append(s, "-Djava.library.path=/usr/share/cassandra/lib/sigar-bin")
	s = append(s, "-Dcassandra.libjemalloc=/usr/lib/x86_64-linux-gnu/libjemalloc.so.1")
	s = append(s, "-Dlogback.configurationFile=file:///etc/cassandra/logback.xml")
	s = append(s, "-Djava.net.preferIPv4Stack=true")
	s = append(s, "-Dcassandra.logdir=/var/log/cassandra")
	s = append(s, "-Dcassandra.config=file:///etc/cassandra/cassandra.yaml")
	s = append(s, "-Dcassandra.storagedir=/var/lib/cassandra/data")
	s = append(s, "-Dcassandra-foreground=yes")
	s = append(s, "-cp")
	s = append(s, "/etc/cassandra:/usr/share/cassandra/apache-cassandra.jar:/usr/share/cassandra/lib/*")
	return s
}
