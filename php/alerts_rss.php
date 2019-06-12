<?php
include("dbconfig.php");

$meta_title = "Wainwrights On The Air - Latest Alerts";
$meta_desc = "Pending activations from Wainwright summits - see http://www.wota.org.uk for more information.";
$meta_link = "http://www.wota.org.uk/alerts_rss.php";

function sanitize_string($string) {
	return htmlentities(strip_tags($string));
}

header("Content-Type: text/xml; charset=iso-8859-1");
print '<?xml version="1.0" encoding="iso-8859-1" ?>' . "\n";
print "<rss version=\"2.0\"  xmlns:atom=\"http://www.w3.org/2005/Atom\">\n";

print "<channel>\n";
print "<title>" . sanitize_string($meta_title) . "</title>\n";
print "<link>" .$meta_link . "</link>\n";
print "<description><![CDATA[" . sanitize_string($meta_desc) . "]]></description>\n";
print "<language>en</language>\n";
print "<atom:link href=\"" . $meta_link . "\" rel=\"self\" type=\"application/rss+xml\" />\n";

$result = mysql_query("SELECT * FROM `alerts` WHERE `datetime` >= CURRENT_DATE ORDER BY `datetime` DESC",$con);
if (mysql_num_rows($result) ) {
  while($row = mysql_fetch_array($result, MYSQL_ASSOC)) {
    $wotaid = $row['wotaid'];
    $res2 = mysql_query("SELECT `name` FROM `summits` WHERE `wotaid` = '".$wotaid."'",$con);
    $row2 = mysql_fetch_array($res2, MYSQL_ASSOC);

 // $wotaid = "LDW-".substr("000".$wotaid,-3,3);

$fellid = (int)$wotaid;
   if ($fellid <= 214) {$prefix = "LDW-";}
  else {$prefix = "LDO-"; 
   $fellid = (int)$fellid-214;}
   $ldid = $prefix . str_pad($fellid,3,"0",STR_PAD_LEFT);
   $wotaid = $ldid;

    $name = $row2['name'];
    $title = $row['call']." on ".$wotaid." - ".$name;
    $desc = "Frequencies/modes: ".$row['freqmode'].". ".$row['comment'].". Posted by ".$row['postedby'].".";
    $date = substr(date("r", strtotime($row['datetime'])),0,26)."+0000";
    //$guiddate = str_replace(" ", "%20", $row['datetime']);

    print "<item>\n";
    print "<title>" . sanitize_string($title) . "</title>\n";
    //print "<guid>http://www.wota.org.uk/alerts/" . $guiddate . "/" . $wotaid . "/" . $row['call'] . "</guid>\n";
    print "<link>http://www.wota.org.uk/alerts/" . $row['id'] . "</link>\n";
    print "<guid>http://www.wota.org.uk/alerts/" . $row['id'] . "</guid>\n";
    print "<description><![CDATA[" . $desc . "]]></description>\n";
  
    print "<pubDate>" . $date . "</pubDate>\n";
    
    print "</item>\n";
  }
} else {
  print "<item>\n";
  print "<title>No alerts</title>\n";
  print "<link>http://www.wota.org.uk/alerts/0</link>\n";
  print "<guid>http://www.wota.org.uk/alerts/0</guid>\n";
  print "</item>\n";
}
print "</channel>\n";
print "</rss>\n";
?>
